import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function ResetPasswordPage() {
  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Reset password</CardTitle>
        <CardDescription>Placeholder route — accept token from search params and post to your API.</CardDescription>
      </CardHeader>
      <CardFooter>
        <Button variant="outline" asChild>
          <Link href="/auth/login">Back to sign in</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
