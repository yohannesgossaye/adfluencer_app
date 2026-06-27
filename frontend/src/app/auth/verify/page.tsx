import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function VerifyPage() {
  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Verify</CardTitle>
        <CardDescription>
          Placeholder route — confirm email or magic-link tokens here.
        </CardDescription>
      </CardHeader>
      <CardFooter>
        <Button variant="outline" asChild>
          <Link href="/auth/login">Go to sign in</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
