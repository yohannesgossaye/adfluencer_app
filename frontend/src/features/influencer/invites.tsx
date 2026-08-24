"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Loader2, MessageSquare } from "lucide-react";

import { Button } from "@/components/ui/button";
import { influencerApi, type InfluencerInviteItem } from "@/lib/api/influencer";
import { chatApi } from "@/lib/api/chat";
import { cn } from "@/lib/utils";

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
  const [isLoading, setIsLoading] = useState(true);
  const [startingChatId, setStartingChatId] = useState<string | null>(null);

  const loadInvites = async (statusFilter?: string) => {
    try {
      const res = await influencerApi.listInvites(
        statusFilter && statusFilter !== "all" ? { status: statusFilter } : undefined
      );
      if (res.data?.items) {
        setInvites(res.data.items);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadInvites(tab);
  }, [tab]);

  const updateStatus = async (id: string, action: "accepted" | "rejected") => {
    try {
      if (action === "accepted") {
        await influencerApi.acceptInvite(id);
      } else {
        await influencerApi.rejectInvite(id);
      }
      await loadInvites(tab);
    } catch {
      // fallback
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
    } catch {
      router.push("/influencer/messages");
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
            onClick={() => setTab(t.id)}
            className={cn(
              "rounded-full border px-4 py-2 text-sm font-medium transition-colors",
              tab === t.id
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:border-white/20 hover:text-white",
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      {invites.length === 0 && !isLoading ? (
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
                    className="h-10 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                    onClick={() => updateStatus(inv.id, "accepted")}
                  >
                    Accept
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    className="h-10 flex-1 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                    onClick={() => updateStatus(inv.id, "rejected")}
                  >
                    Reject
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
