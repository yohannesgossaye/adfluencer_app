import { ApprovalGate } from "@/components/influencer/account-status";
import { InfluencerDashboard } from "@/features/influencer/dashboard";

export default function InfluencerDashboardPage() {
  return (
    <ApprovalGate feature="Invites, contracts and messages">
      <InfluencerDashboard />
    </ApprovalGate>
  );
}
