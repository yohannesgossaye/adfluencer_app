import Link from "next/link";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function ForgotPasswordPage() {
  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <CardTitle>Forgot password</CardTitle>
        <CardDescription>Placeholder route — wire your reset flow here.</CardDescription>
      </CardHeader>
      <CardFooter>
        <Button variant="outline" asChild>
          <Link href="/auth/login">Back to sign in</Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
