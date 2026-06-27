import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Influencer } from "@/types";

import { StatusBadge } from "@/components/common/status-badge";

export function InfluencerCard({ influencer }: { influencer: Influencer }) {
  return (
    <Card className="adfluence-glass border border-white/10 text-white">
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <CardTitle className="text-base">{influencer.name}</CardTitle>
        <StatusBadge status={influencer.status} />
      </CardHeader>
      <CardContent className="space-y-1 text-sm text-[#CBD5E1]">
        <p className="capitalize">{influencer.niche}</p>
        <p>{influencer.city}</p>
        <p>{influencer.followers.toLocaleString()} followers</p>
      </CardContent>
    </Card>
  );
}
