import { StyledAuthLink } from "@/components/features/auth/AuthCompoents";
import LoginForm from "@/components/features/auth/LoginForm";

export default function LoginPage() {
  return (
    <>
      <LoginForm />
      <StyledAuthLink href={"/register"}>
        Don`&apos;`t have an account? Create one
      </StyledAuthLink>
    </>
  );
}
