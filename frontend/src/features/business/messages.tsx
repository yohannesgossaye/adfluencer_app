"use client";

import Link from "next/link";

import { Button } from "@/components/ui/button";
import { MessagesPage, type MessagesPageConfig } from "@/features/chat/messages-page";

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
      <Button
        asChild
        size="sm"
        className="h-9 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:flex-none"
      >
        <Link href="/business/contracts">Create contract</Link>
      </Button>
    </>
  ),
};

export function BusinessMessages() {
  return <MessagesPage config={config} />;
}
