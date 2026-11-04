/** Support conversation types shared by the business, creator and admin portals (backend dto/support). */

export type SupportMessage = {
  id: string;
  from: "user" | "support";
  sender_name: string;
  body: string;
  created_at: string;
};

/** Who wrote last: "support" (our team), "user", or "" before anyone replied. */
export type LastReplyBy = "" | "user" | "support";

export type SupportReplyResult = {
  message: SupportMessage;
  status: string;
  reopened: boolean;
};

export const SUPPORT_MESSAGE_MAX = 4000;

export function isTicketDone(status: string) {
  return status === "resolved" || status === "closed";
}
