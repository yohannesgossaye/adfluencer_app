import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Contract } from "@/types";

import { StatusBadge } from "@/components/common/status-badge";

export function ContractCard({ contract }: { contract: Contract }) {
  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <CardTitle className="text-base">{contract.campaignTitle}</CardTitle>
        <StatusBadge status={contract.status} />
      </CardHeader>
      <CardContent className="space-y-1 text-sm text-muted-foreground">
        <p>{contract.businessName}</p>
        <p>{contract.influencerName}</p>
        <p>ETB {contract.amountEtb.toLocaleString()}</p>
      </CardContent>
    </Card>
  );
}
