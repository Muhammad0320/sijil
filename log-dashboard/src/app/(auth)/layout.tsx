import {
  AuthCard,
  AuthContainer,
  AuthTitle,
} from "@/components/features/auth/AuthCompoents";
import { ReactNode } from "react";

export default function AuthLayout({ children }: { children: ReactNode }) {
  return (
    <AuthContainer>
      <AuthCard>
        <AuthTitle>LogEngine</AuthTitle>
        {children}
      </AuthCard>
    </AuthContainer>
  );
}
