import { ApprovalGate } from "@/components/influencer/account-status";
import { InfluencerContracts } from "@/features/influencer/contracts";

export default function InfluencerContractsPage() {
  return (
    <ApprovalGate feature="Contracts">
      <InfluencerContracts />
    </ApprovalGate>
  );
}
