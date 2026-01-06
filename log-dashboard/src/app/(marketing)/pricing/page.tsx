import PricingSection from "@/components/marketing/PricingSection";
import styled from "styled-components";

const Container = styled.div`
  padding-top: 100px;
  min-height: 100vh;
  background: #050505;
`;

export default function PricingPage() {
  return (
    <Container>
      <PricingSection />
    </Container>
  );
}
