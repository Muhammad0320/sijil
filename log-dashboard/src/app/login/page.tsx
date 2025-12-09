import LoginForm from "@/components/features/auth/LoginForm";
import { LoginCard, LoginContainer, LoginTitle } from "./LoginLayout";

export default function LoginPage() {
  return (
    <LoginContainer>
      <LoginCard>
        <LoginTitle>LogEngine Access</LoginTitle>
        {/* The interactive island */}
        <LoginForm />
      </LoginCard>
    </LoginContainer>
  );
}
