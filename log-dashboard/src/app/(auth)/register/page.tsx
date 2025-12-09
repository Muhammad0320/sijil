import { AuthLink } from "@/components/features/auth/AuthCompoents";
import RegisterForm from "@/components/features/auth/RegisterForm";

import Link from "next/link";

export default function RegisterPage() {
  return (
    <>
      <RegisterForm />
      <Link href="/login" passHref legacyBehavior>
        <AuthLink>Already have an account? Sign in</AuthLink>
      </Link>
    </>
  );
}
