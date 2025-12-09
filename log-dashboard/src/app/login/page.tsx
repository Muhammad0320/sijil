import LoginForm from "@/components/features/auth/LoginForm";
import styled from "styled-components";

// Reuse these styled components or move them to a shared layout file
const Container = styled.div`
  height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-color);
`;

const Card = styled.div`
  width: 400px;
  padding: 40px;
  background: #161b22;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
`;

const Title = styled.h1`
  color: var(--text-color);
  font-size: 24px;
  text-align: center;
  margin-bottom: 24px;
`;

export default function LoginPage() {
  return (
    <Container>
      <Card>
        <Title>LogEngine Access</Title>
        {/* The interactive island */}
        <LoginForm />
      </Card>
    </Container>
  );
}
