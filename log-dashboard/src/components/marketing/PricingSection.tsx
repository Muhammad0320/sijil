import { Check, Info } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import styled, { css } from "styled-components";

// --- PRICING ---
const PricingSectionContainer = styled.section`
  padding: 100px 24px;
  max-width: 1200px;
  margin: 0 auto;
`;

const ToggleContainer = styled.div`
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 64px;
`;

const ToggleBtn = styled.button<{ $active: boolean }>`
  padding: 10px 24px;
  border-radius: 99px;
  background: ${(p) => (p.$active ? "#238636" : "rgba(255,255,255,0.05)")};
  color: #fff;
  font-weight: 600;
  border: 1px solid ${(p) => (p.$active ? "#238636" : "#30363d")};
  transition: all 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94);

  &:hover {
    background: ${(p) => (p.$active ? "#2ea043" : "rgba(255,255,255,0.1)")};
  }
`;

const PricingGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  align-items: stretch; /* Fix height alignment */
`;

const CardTitle = styled.h3`
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
`;

const PrimaryButton = styled(Link)`
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  color: #000;
  padding: 16px 32px;
  border-radius: 8px;
  font-weight: 700;
  font-size: 16px;
  transition: transform 0.2s;
  &:hover {
    transform: translateY(-2px);
  }
`;

const PriceCard = styled.div<{ $featured?: boolean; $ultra?: boolean }>`
  background: ${(p) => (p.$featured ? "rgba(88, 166, 255, 0.05)" : "#161b22")};
  border: 1px solid ${(p) => (p.$featured ? "#58a6ff" : "#30363d")};
  border-radius: 24px;
  padding: 40px 32px;
  display: flex;
  flex-direction: column;
  position: relative;
  transition: transform 0.3s;
  height: 100%;

  &:hover {
    transform: translateY(-8px);
  }

  ${(p) =>
    p.$ultra &&
    css`
      background: linear-gradient(#0d1117, #0d1117) padding-box,
        linear-gradient(45deg, #ffd700, #ffaa00) border-box;
      border: 1px solid transparent;
    `}
`;

const PriceAmount = styled.div`
  font-size: 42px;
  font-weight: 800;
  color: #fff;
  margin: 24px 0 8px;
  font-family: var(--font-geist-mono);
  display: flex;
  align-items: baseline;
  gap: 4px;
  span {
    font-size: 16px;
    color: #8b949e;
    font-weight: 400;
  }
`;

const BilledText = styled.div`
  font-size: 13px;
  color: #2ecc71;
  font-weight: 600;
  height: 20px;
`;

const FeatureList = styled.ul`
  list-style: none;
  margin: 32px 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1; /* Pushes button to bottom */
`;

const FeatureItem = styled.li`
  display: flex;
  align-items: center;
  gap: 12px;
  color: #c9d1d9;
  font-size: 14px;
`;

const SectionHeader = styled.div`
  text-align: center;
  margin-bottom: 64px;
  h2 {
    font-size: 48px;
    font-weight: 800;
    color: #fff;
    margin-bottom: 16px;
  }
  p {
    font-size: 18px;
    color: #8b949e;
  }
`;

export default function PricingSection() {
  const [currency, setCurrency] = useState<"USD" | "NGN">("USD");
  const [isAnnual, setIsAnnual] = useState(true);

  const rates = {
    usd: { pro: 20, ultra: 100 },
    ngn: { pro: 12500, ultra: 95000 },
  };
  const currentRate = currency === "USD" ? rates.usd : rates.ngn;
  const symbol = currency === "USD" ? "$" : "₦";

  // Calculate Monthly Equivalent
  const getPrice = (base: number) => {
    if (!isAnnual) return base.toLocaleString();
    const annualTotal = base * 10; // 2 Months Free
    return (annualTotal / 12).toLocaleString(undefined, {
      maximumFractionDigits: 0,
    });
  };

  const getBilledText = (base: number) => {
    if (!isAnnual) return "Billed monthly";
    const annualTotal = base * 10;
    return `Billed ${symbol}${annualTotal.toLocaleString()} yearly (Save 17%)`;
  };

  return (
    <PricingSectionContainer id="pricing">
      <SectionHeader>
        <h2>Transparent Pricing</h2>
        <p>Start small, scale infinitely.</p>
      </SectionHeader>

      <ToggleContainer>
        <div
          style={{
            display: "flex",
            gap: 4,
            background: "#161b22",
            padding: 4,
            borderRadius: 99,
            border: "1px solid #30363d",
          }}
        >
          <ToggleBtn
            $active={currency === "USD"}
            onClick={() => setCurrency("USD")}
          >
            🇺🇸 USD
          </ToggleBtn>
          <ToggleBtn
            $active={currency === "NGN"}
            onClick={() => setCurrency("NGN")}
          >
            🇳🇬 NGN
          </ToggleBtn>
        </div>
        <div
          style={{
            display: "flex",
            gap: 4,
            background: "#161b22",
            padding: 4,
            borderRadius: 99,
            border: "1px solid #30363d",
          }}
        >
          <ToggleBtn $active={!isAnnual} onClick={() => setIsAnnual(false)}>
            Monthly
          </ToggleBtn>
          <ToggleBtn $active={isAnnual} onClick={() => setIsAnnual(true)}>
            Yearly{" "}
            <span style={{ color: "#4ade80", fontSize: 11 }}>(-17%)</span>
          </ToggleBtn>
        </div>
      </ToggleContainer>

      <PricingGrid>
        {/* Hobby */}
        <PriceCard>
          <CardTitle>Hobby</CardTitle>
          <PriceAmount>
            $0<span>/mo</span>
          </PriceAmount>
          <BilledText>Free forever</BilledText>
          <FeatureList>
            <FeatureItem>
              <Check size={16} color="#2ecc71" /> 10k Logs/Day
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#2ecc71" /> 3 Day Retention
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#2ecc71" /> 1 Member
            </FeatureItem>
          </FeatureList>
          <PrimaryButton
            href="/register"
            style={{
              justifyContent: "center",
              width: "100%",
              background: "rgba(255,255,255,0.1)",
            }}
          >
            Start Free
          </PrimaryButton>
        </PriceCard>

        {/* Pro */}
        <PriceCard $featured>
          <div
            style={{
              position: "absolute",
              top: -12,
              left: "50%",
              transform: "translateX(-50%)",
              background: "#58a6ff",
              color: "#000",
              fontSize: 10,
              fontWeight: 800,
              padding: "4px 12px",
              borderRadius: 20,
            }}
          >
            RECOMMENDED
          </div>
          <CardTitle>Pro</CardTitle>
          <PriceAmount>
            {symbol}
            {getPrice(currentRate.pro)}
            <span>/mo</span>
          </PriceAmount>
          <BilledText>{getBilledText(currentRate.pro)}</BilledText>
          <FeatureList>
            <FeatureItem>
              <Check size={16} color="#58a6ff" /> 1M Logs/Day
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#58a6ff" /> 14 Day Retention
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#58a6ff" /> 10 Projects
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#58a6ff" /> Unlimited Members
            </FeatureItem>
          </FeatureList>
          <PrimaryButton
            href="/register"
            style={{
              justifyContent: "center",
              width: "100%",
              background: "#58a6ff",
              color: "#000",
            }}
          >
            Get Pro
          </PrimaryButton>
        </PriceCard>

        {/* Ultra */}
        <PriceCard $ultra>
          <CardTitle
            style={{
              color: "#ffd700",
              display: "flex",
              alignItems: "center",
              gap: 8,
            }}
          >
            Ultra <Info size={16} />
          </CardTitle>
          <PriceAmount>
            {symbol}
            {getPrice(currentRate.ultra)}
            <span>/mo</span>
          </PriceAmount>
          <BilledText>{getBilledText(currentRate.ultra)}</BilledText>
          <FeatureList>
            <FeatureItem>
              <Check size={16} color="#ffd700" /> 20M Logs/Day
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#ffd700" /> 30 Day Retention
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#ffd700" /> Unlimited Projects
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#ffd700" /> Unlimited Members
            </FeatureItem>
            <FeatureItem>
              <Check size={16} color="#ffd700" /> Cold Storage (S3)
            </FeatureItem>
          </FeatureList>
          <PrimaryButton
            href="mailto:sales@sijil.dev"
            style={{
              justifyContent: "center",
              width: "100%",
              background: "linear-gradient(to right, #ffd700, #ffaa00)",
              color: "#000",
            }}
          >
            Contact Sales
          </PrimaryButton>
        </PriceCard>
      </PricingGrid>
    </PricingSectionContainer>
  );
}
