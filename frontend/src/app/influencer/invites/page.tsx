import { ApprovalGate } from "@/components/influencer/account-status";
import { InfluencerInvites } from "@/features/influencer/invites";

export default function InfluencerInvitesPage() {
  return (
    <ApprovalGate feature="Campaign invites">
      <InfluencerInvites />
    </ApprovalGate>
  );
}
