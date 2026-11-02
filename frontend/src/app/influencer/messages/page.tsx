import { ApprovalGate } from "@/components/influencer/account-status";
import { InfluencerMessages } from "@/features/influencer/messages";

export default function InfluencerMessagesPage() {
  return (
    <ApprovalGate feature="Messages">
      <InfluencerMessages />
    </ApprovalGate>
  );
}
