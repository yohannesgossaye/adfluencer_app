"use client";

import type { ColumnDef } from "@tanstack/react-table";

import type { Organization } from "@/lib/api/types";

export const organizationColumns: ColumnDef<Organization>[] = [
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "slug",
    header: "Slug",
  },
  {
    accessorKey: "createdAt",
    header: "Created",
    cell: ({ row }) =>
      new Date(row.original.createdAt).toLocaleDateString(undefined, {
        dateStyle: "medium",
      }),
  },
];
