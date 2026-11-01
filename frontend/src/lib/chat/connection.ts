import { CHAT_WS_SUBPROTOCOL, type WSIncomingMessage, type WSOutgoingMessage } from "@/lib/api/chat";

export type ChatConnectionState = "idle" | "connecting" | "open" | "reconnecting" | "stopped";

type Options = {
  url: string;
  /** Returns the current access token. `fresh` asks for a refreshed one. */
  getToken: (fresh: boolean) => Promise<string | null>;
  onEvent: (event: WSOutgoingMessage) => void;
  onStateChange: (state: ChatConnectionState) => void;
  /** Called after a reconnect so the UI can reload anything it missed. */
  onResync: () => void;
};

// Close codes sent by the backend (see realtime/client.go).
const CLOSE_TOKEN_EXPIRED = 4001;
const CLOSE_ACCOUNT_UNAVAILABLE = 4003;

const MIN_BACKOFF_MS = 1_000;
const MAX_BACKOFF_MS = 30_000;

/**
 * One realtime chat connection per signed-in session. Framework-free so the
 * web portal and the React Native app can share it.
 *
 * - Authenticates with the access token as a WebSocket subprotocol.
 * - Reconnects with exponential backoff + jitter; refreshes the token first
 *   when the server says it expired (4001) or the handshake was refused.
 * - Stops for good when the account may no longer chat (4003) or on stop().
 */
export class ChatConnection {
  private ws: WebSocket | null = null;
  private attempts = 0;
  private retryTimer: ReturnType<typeof setTimeout> | null = null;
  private stopped = false;
  private needsFreshToken = false;
  private hasConnectedBefore = false;
  private state: ChatConnectionState = "idle";

  constructor(private readonly opts: Options) {}

  get currentState() {
    return this.state;
  }

  start() {
    this.stopped = false;
    void this.connect();
  }

  stop() {
    this.stopped = true;
    this.clearRetry();
    const ws = this.ws;
    this.ws = null;
    ws?.close(1000, "client closed");
    this.setState("stopped");
  }

  /** Reconnect right away (e.g. tab became visible or network came back). */
  nudge() {
    if (this.stopped || this.state === "open" || this.state === "connecting") return;
    this.clearRetry();
    void this.connect();
  }

  /** Sends a frame if the socket is open. Returns false otherwise. */
  send(frame: WSIncomingMessage): boolean {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return false;
    this.ws.send(JSON.stringify(frame));
    return true;
  }

  private setState(next: ChatConnectionState) {
    if (this.state === next) return;
    this.state = next;
    this.opts.onStateChange(next);
  }

  private clearRetry() {
    if (this.retryTimer) {
      clearTimeout(this.retryTimer);
      this.retryTimer = null;
    }
  }

  private async connect() {
    if (this.stopped) return;
    this.setState(this.hasConnectedBefore ? "reconnecting" : "connecting");

    let token: string | null = null;
    try {
      token = await this.opts.getToken(this.needsFreshToken);
    } catch {
      token = null;
    }
    if (this.stopped) return;
    if (!token) {
      this.scheduleRetry();
      return;
    }

    let opened = false;
    let ws: WebSocket;
    try {
      ws = new WebSocket(this.opts.url, [CHAT_WS_SUBPROTOCOL, token]);
    } catch {
      this.scheduleRetry();
      return;
    }
    this.ws = ws;

    ws.onopen = () => {
      opened = true;
      this.attempts = 0;
      this.needsFreshToken = false;
      const isReconnect = this.hasConnectedBefore;
      this.hasConnectedBefore = true;
      this.setState("open");
      if (isReconnect) this.opts.onResync();
    };

    ws.onmessage = (msg) => {
      try {
        this.opts.onEvent(JSON.parse(String(msg.data)) as WSOutgoingMessage);
      } catch {
        // Ignore frames that are not valid JSON.
      }
    };

    ws.onclose = (event) => {
      if (this.ws !== ws) return; // superseded by stop() or a newer socket
      this.ws = null;
      if (this.stopped) return;
      if (event.code === CLOSE_ACCOUNT_UNAVAILABLE) {
        this.setState("stopped");
        return;
      }
      // A refused handshake usually means the access token expired.
      this.needsFreshToken = event.code === CLOSE_TOKEN_EXPIRED || !opened;
      this.scheduleRetry(event.code === CLOSE_TOKEN_EXPIRED);
    };
  }

  private scheduleRetry(immediate = false) {
    if (this.stopped) return;
    this.setState("reconnecting");
    this.clearRetry();
    const base = Math.min(MAX_BACKOFF_MS, MIN_BACKOFF_MS * 2 ** this.attempts);
    const delay = immediate ? 0 : base / 2 + Math.random() * (base / 2);
    this.attempts += 1;
    this.retryTimer = setTimeout(() => void this.connect(), delay);
  }
}
