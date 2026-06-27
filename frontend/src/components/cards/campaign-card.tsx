import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Campaign } from "@/types";

import { StatusBadge } from "@/components/common/status-badge";

export function CampaignCard({ campaign }: { campaign: Campaign }) {
  return (
    <Card className="adfluence-glass border border-white/10 text-white">
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <CardTitle className="text-base">{campaign.title}</CardTitle>
        <StatusBadge status={campaign.status} />
      </CardHeader>
      <CardContent className="space-y-1 text-sm text-[#CBD5E1]">
        <p>{campaign.businessName}</p>
        <p>{campaign.platform}</p>
        <p>ETB {campaign.budgetEtb.toLocaleString()}</p>
      </CardContent>
    </Card>
  );
}
