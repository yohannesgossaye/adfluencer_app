"use client";

import { MessagesPage, type MessagesPageConfig } from "@/features/chat/messages-page";

const config: MessagesPageConfig = {
  subtitle: "Brands & campaigns",
  searchPlaceholder: "Search",
  emptyHint: "When you accept a campaign invite or a brand messages you, chats will appear here.",
  noSelection: "Select a conversation to start chatting in real time.",
  quickReplies: ["Yes, I’m interested", "Can we discuss the budget?", "Please send the contract"],
};

export function InfluencerMessages() {
  return <MessagesPage config={config} />;
}
