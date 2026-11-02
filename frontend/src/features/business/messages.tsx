"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { MessagesPage, type MessagesPageConfig } from "@/features/chat/messages-page";
import { businessApi } from "@/lib/api/business";
import type { ChatConversationItem } from "@/lib/api/chat";

const primary =
  "h-9 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:flex-none";

/**
 * A contract needs an accepted invite. When this chat belongs to a campaign the
 * creator accepted, open the contract form already filled in; otherwise point
 * to inviting them first instead of an empty form.
 */
function ContractAction({ conversation }: { conversation: ChatConversationItem }) {
  const [inviteId, setInviteId] = useState<string | null>(null);
  const campaignId = conversation.campaign_id;
  const creatorId = conversation.other_user_id;

  useEffect(() => {
    if (!campaignId) return;
    let cancelled = false;
    businessApi
      .listInvites(campaignId)
      .then((res) => {
        const accepted = res.data?.items?.find((i) => i.influencer_user_id === creatorId && i.status === "accepted");
        if (!cancelled) setInviteId(accepted?.id ?? null);
      })
      .catch(() => !cancelled && setInviteId(null));
    return () => {
      cancelled = true;
    };
  }, [campaignId, creatorId]);

  if (campaignId && inviteId) {
    return (
      <Button asChild size="sm" className={primary}>
        <Link href={`/business/contracts?campaign_id=${campaignId}&invite_id=${inviteId}`}>Send contract</Link>
      </Button>
    );
  }
  return (
    <Button asChild size="sm" className={primary}>
      <Link href={`/business/influencers/${creatorId}`}>Invite to a campaign</Link>
    </Button>
  );
}

const config: MessagesPageConfig = {
  subtitle: "Campaign conversations",
  searchPlaceholder: "Search conversations",
  emptyHint: "When you invite a creator or message them, chats will appear here.",
  noSelection: "Select a conversation to start chatting with creators in real time.",
  quickReplies: ["Can you share your rate?", "Let’s confirm deliverables", "Can we start this week?"],
  headerActions: (conversation) => (
    <>
      <Button asChild size="sm" className="h-9 flex-1 rounded-xl bg-white/10 text-white hover:bg-white/15 sm:flex-none">
        <Link href={`/business/influencers/${conversation.other_user_id}`}>View profile</Link>
      </Button>
      <ContractAction key={conversation.id} conversation={conversation} />
    </>
  ),
};

export function BusinessMessages() {
  return <MessagesPage config={config} />;
}
