import { AuthLink } from "@/components/features/auth/AuthCompoents";
import LoginForm from "@/components/features/auth/LoginForm";
import Link from "next/link";

export default function LoginPage() {
  return (
    <>
      <LoginForm />
      <Link href="/register" passHref legacyBehavior>
        <AuthLink>Don`&apos;`t have an account? Create one</AuthLink>
      </Link>
    </>
  );
}
