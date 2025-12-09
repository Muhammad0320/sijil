import { StyledAuthLink } from "@/components/features/auth/AuthCompoents";
import RegisterForm from "@/components/features/auth/RegisterForm";

export default function RegisterPage() {
  return (
    <>
      <RegisterForm />
      <StyledAuthLink href={"/login"}>
        Already have an account? sign in.
      </StyledAuthLink>
    </>
  );
}
