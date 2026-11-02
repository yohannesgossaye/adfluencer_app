"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Loader2, MessageSquare } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { influencerApi, type InfluencerInviteItem } from "@/lib/api/influencer";
import { chatApi } from "@/lib/api/chat";
import { cn } from "@/lib/utils";

function apiMessage(err: unknown, fallback: string) {
  const data = (err as { response?: { data?: { message?: string; error?: { type?: string } } } })?.response?.data;
  if (data?.error?.type === "invalid_campaign_status") return "This campaign is no longer running, so the invite can't be changed.";
  if (data?.error?.type === "account_not_verified") return "Your profile must be approved before you can accept invites.";
  return data?.message || fallback;
}

const tabs: Array<{ id: "all" | "pending" | "accepted" | "rejected"; label: string }> = [
  { id: "all", label: "All" },
  { id: "pending", label: "Pending" },
  { id: "accepted", label: "Accepted" },
  { id: "rejected", label: "Rejected" },
];

function statusStyle(status: string) {
  const s = status.toLowerCase();
  switch (s) {
    case "pending":
      return "border-amber-400/30 bg-amber-500/10 text-amber-200";
    case "accepted":
      return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
    case "rejected":
      return "border-rose-400/30 bg-rose-500/10 text-rose-200";
    default:
      return "border-white/10 bg-white/5 text-[#CBD5E1]";
  }
}

export function InfluencerInvites() {
  const router = useRouter();
  const [invites, setInvites] = useState<InfluencerInviteItem[]>([]);
  const [tab, setTab] = useState<(typeof tabs)[number]["id"]>("all");
  const [listState, setListState] = useState<"loading" | "ready" | "error">("loading");
  const [busyId, setBusyId] = useState<string | null>(null);
  const [startingChatId, setStartingChatId] = useState<string | null>(null);

  const loadInvites = async (statusFilter?: string) => {
    try {
      const res = await influencerApi.listInvites(
        statusFilter && statusFilter !== "all" ? { status: statusFilter } : undefined
      );
      setInvites(res.data?.items ?? []);
      setListState("ready");
    } catch {
      setListState("error");
    }
  };

  useEffect(() => {
    void Promise.resolve().then(() => loadInvites(tab));
  }, [tab]);

  const updateStatus = async (id: string, action: "accepted" | "rejected") => {
    if (action === "rejected" && !window.confirm("Decline this invite? The business will be told.")) return;
    setBusyId(id);
    try {
      if (action === "accepted") {
        await influencerApi.acceptInvite(id);
      } else {
        await influencerApi.rejectInvite(id);
      }
      toast.success(action === "accepted" ? "Invite accepted. The business can now send you a contract." : "Invite declined.");
      await loadInvites(tab);
    } catch (err) {
      toast.error(apiMessage(err, "That didn't work. Please try again."));
    } finally {
      setBusyId(null);
    }
  };

  const handleOpenChat = async (inv: InfluencerInviteItem) => {
    const businessUserId = inv.business_user_id || inv.campaign?.business_user_id;
    setStartingChatId(inv.id);
    try {
      if (businessUserId) {
        const res = await chatApi.getOrCreateConversation({
          recipient_user_id: businessUserId,
          campaign_id: inv.campaign?.id,
        });
        if (res.data?.id) {
          router.push(`/influencer/messages?conversation_id=${res.data.id}`);
          return;
        }
      }
      router.push("/influencer/messages");
    } catch (err) {
      toast.error(apiMessage(err, "Chat could not be opened. Please try again."));
    } finally {
      setStartingChatId(null);
    }
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Campaign invites</h1>
        <p className="mt-1 text-sm text-[#CBD5E1]">
          Review invitations from brands and respond on your timeline.
        </p>
      </div>

      <div className="flex flex-wrap gap-2">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => { setTab(t.id); setListState("loading"); }}
            className={cn(
              "rounded-full border px-4 py-2 text-sm font-medium transition-colors",
              tab === t.id
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:border-white/20 hover:text-white",
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      {listState === "loading" ? (
        <p role="status" className="flex items-center gap-2 p-6 text-sm text-[#94A3B8]">
          <Loader2 className="size-4 animate-spin" aria-hidden /> Loading invites…
        </p>
      ) : listState === "error" ? (
        <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
          Invites could not be loaded.{" "}
          <button type="button" className="underline" onClick={() => { setListState("loading"); void loadInvites(tab); }}>
            Try again
          </button>
        </div>
      ) : invites.length === 0 ? (
        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
          <p className="text-white">No campaign invites found.</p>
          <p className="mt-1 text-sm text-[#CBD5E1]">When businesses invite you to their campaigns, they will appear here.</p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {invites.map((inv) => (
            <article
              key={inv.id}
              className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <p className="font-semibold text-white">{inv.campaign?.business_name || "Brand"}</p>
                  <p className="text-sm text-[#CBD5E1]">{inv.campaign?.title}</p>
                </div>
                <span
                  className={cn("rounded-full border px-2.5 py-0.5 text-xs capitalize", statusStyle(inv.status))}
                >
                  {inv.status}
                </span>
              </div>
              <dl className="mt-3 space-y-1 text-sm text-[#CBD5E1]">
                <div>
                  <dt className="inline text-[#94A3B8]">Platforms: </dt>
                  <dd className="inline text-white capitalize">{inv.campaign?.platforms?.join(", ") || "Multi"}</dd>
                </div>
                <div>
                  <dt className="inline text-[#94A3B8]">Budget: </dt>
                  <dd className="inline text-white">
                    {inv.campaign?.budget_min_etb != null && inv.campaign?.budget_max_etb != null
                      ? `ETB ${inv.campaign.budget_min_etb.toLocaleString()} - ${inv.campaign.budget_max_etb.toLocaleString()}`
                      : "Flexible"}
                  </dd>
                </div>
                {inv.campaign?.start_date && inv.campaign?.end_date ? (
                  <div>
                    <dt className="inline text-[#94A3B8]">Timeline: </dt>
                    <dd className="inline text-white">
                      {new Date(inv.campaign.start_date).toLocaleDateString()} - {new Date(inv.campaign.end_date).toLocaleDateString()}
                    </dd>
                  </div>
                ) : null}
              </dl>
              {inv.message ? (
                <p className="mt-3 rounded-xl border border-white/10 bg-[#06111F] p-3 text-sm text-[#CBD5E1]">
                  {inv.message}
                </p>
              ) : null}
              {inv.status === "pending" ? (
                <div className="mt-4 flex flex-col gap-2 sm:flex-row">
                  <Button
                    type="button"
                    className="h-10 flex-1 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]"
                    disabled={busyId === inv.id}
                    onClick={() => updateStatus(inv.id, "accepted")}
                  >
                    Accept
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    className="h-10 flex-1 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                    disabled={busyId === inv.id}
                    onClick={() => updateStatus(inv.id, "rejected")}
                  >
                    Decline
                  </Button>
                </div>
              ) : inv.status === "accepted" ? (
                <div className="mt-4 flex flex-col gap-2 sm:flex-row">
                  <Button
                    type="button"
                    disabled={startingChatId === inv.id}
                    className="h-10 flex-1 rounded-xl bg-emerald-600 text-white hover:bg-emerald-500 flex items-center justify-center gap-2"
                    onClick={() => handleOpenChat(inv)}
                  >
                    {startingChatId === inv.id ? (
                      <>
                        <Loader2 className="size-4 animate-spin" />
                        Opening Chat...
                      </>
                    ) : (
                      <>
                        <MessageSquare className="size-4" />
                        Chat with {inv.campaign?.business_name || "Brand"}
                      </>
                    )}
                  </Button>
                </div>
              ) : null}
            </article>
          ))}
        </div>
      )}
    </div>
  );
}
