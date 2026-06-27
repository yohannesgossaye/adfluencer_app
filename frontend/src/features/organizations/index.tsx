"use client";

import { useQuery } from "@tanstack/react-query";
import {
  getCoreRowModel,
  useReactTable,
  type VisibilityState,
} from "@tanstack/react-table";
import * as React from "react";

import { DataTable } from "@/components/table/data-table";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { listOrganizations } from "@/lib/api/organizations";

import { organizationColumns } from "./columns";

export default function Organizations() {
  const [columnVisibility, setColumnVisibility] =
    React.useState<VisibilityState>({});

  const { data = [], isPending } = useQuery({
    queryKey: ["organizations", "list"] as const,
    queryFn: listOrganizations,
  });

  /* eslint-disable-next-line react-hooks/incompatible-library -- TanStack Table */
  const table = useReactTable({
    data,
    columns: organizationColumns,
    state: { columnVisibility },
    onColumnVisibilityChange: setColumnVisibility,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Organizations</h1>
        <p className="text-muted-foreground text-sm">
          Example list wired through the shared API client and React Query.
        </p>
      </div>
      <Card>
        <CardHeader>
          <CardTitle>Directory</CardTitle>
          <CardDescription>
            {isPending ? "Loading…" : `${data.length} organizations`}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <DataTable table={table} emptyMessage="No organizations found." />
        </CardContent>
      </Card>
    </div>
  );
}
