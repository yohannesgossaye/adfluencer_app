import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { StatItem } from "@/types";

export function StatCard({ item }: { item: StatItem }) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm text-slate-500">{item.label}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-2xl font-semibold text-slate-900">{item.value}</p>
        {item.hint ? <p className="text-xs text-slate-500">{item.hint}</p> : null}
      </CardContent>
    </Card>
  );
}
