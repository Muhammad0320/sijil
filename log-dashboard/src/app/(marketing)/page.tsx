"use client";

import React, { useState } from "react";
import styled, { keyframes, css } from "styled-components";
import {
  ArrowRight,
  Database,
  Zap,
  Shield,
  Cpu,
  Check,
  Info,
} from "lucide-react";
import Link from "next/link";
import Image from "next/image";

// Import Assets
import hypercube from "../../../public/hypercube.png";
import vault from "../../../public/vault-green.png";
import server from "../../../public/server.png";
import cacheImg from "../../../public/cache.png";
import websocket from "../../../public/websocket.png";
import HeroRollingText, {
  HandDrawnHighlight,
} from "@/components/marketing/HeroText";
import BenchmarkSection from "@/components/marketing/benchmark";

// --- ANIMATIONS ---
const roll = keyframes`
  0% { transform: translateY(0); }
  25% { transform: translateY(-100%); }
  50% { transform: translateY(-200%); }
  75% { transform: translateY(-300%); }
  100% { transform: translateY(0); }
`;

const draw = keyframes`
  to { stroke-dashoffset: 0; }
`;

const float = keyframes`
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
`;

const scanlineAnim = keyframes`
  0% { top: -20%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
`;

// --- STYLED COMPONENTS ---
const Container = styled.div`
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
  background: #050505;
  color: #c9d1d9;
  font-family: var(--font-geist-sans);
`;

// --- HERO SECTION ---
const HeroSection = styled.section`
  min-height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  padding-top: 40px;

  background-image: linear-gradient(
      rgba(255, 255, 255, 0.03) 1px,
      transparent 1px
    ),
    linear-gradient(90deg, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  background-size: 50px 50px;
  mask-image: radial-gradient(circle at center, black 40%, transparent 80%);
`;

const HeroGlow = styled.div`
  position: absolute;
  width: 600px;
  height: 600px;
  background: radial-gradient(
    circle,
    rgba(88, 166, 255, 0.15) 0%,
    transparent 70%
  );
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  filter: blur(80px);
  z-index: 0;
  pointer-events: none;
`;

const Pill = styled.div`
  background: rgba(88, 166, 255, 0.1);
  border: 1px solid rgba(88, 166, 255, 0.3);
  padding: 6px 16px;
  border-radius: 99px;
  font-size: 13px;
  font-weight: 600;
  color: #58a6ff;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 40px;
  backdrop-filter: blur(4px);
  box-shadow: 0 0 15px rgba(88, 166, 255, 0.2);
`;

// Rolling Text Styles
const RollingContainer = styled.div`
  display: inline-block;
  height: 1.1em;
  overflow: hidden;
  vertical-align: bottom;
  position: relative;
  width: auto;
  min-width: 320px;
  text-align: left;
`;

const RollList = styled.ul`
  margin: 0;
  padding: 0;
  list-style: none;
  animation: ${roll} 8s cubic-bezier(0.23, 1, 0.32, 1) infinite;
`;

const RollItem = styled.li`
  height: 1.1em;
  display: flex;
  align-items: center;
  color: #58a6ff;
  font-weight: 800;
  padding-left: 12px;
`;

const HeroTitle = styled.h1`
  font-size: 72px;
  font-weight: 800;
  letter-spacing: -2px;
  line-height: 1.1;
  text-align: center;
  margin-bottom: 32px;
  z-index: 1;
  max-width: 1200px;
  color: #fff;

  @media (max-width: 768px) {
    font-size: 48px;
  }
`;

const PaintedUnderlineSVG = styled.svg`
  position: absolute;
  bottom: -15px;
  left: 0;
  width: 100%;
  height: 20px;
  z-index: -1;

  path {
    stroke: #58a6ff;
    stroke-width: 8;
    stroke-linecap: round;
    stroke-dasharray: 400;
    stroke-dashoffset: 400;
    fill: none;
    opacity: 0.6;
    animation: ${draw} 1.5s ease-out forwards 0.5s;
  }
`;

const Subtitle = styled.div`
  font-size: 22px;
  color: #8b949e;
  text-align: center;
  max-width: 700px;
  margin-bottom: 48px;
  line-height: 1.6;
  z-index: 1;
`;

const UPS = styled.div`
  font-size: 28px;
  font-weight: 800;
  color: #fff;
  margin-top: 24px;
  text-shadow: 0 0 30px rgba(88, 166, 255, 0.3);
  letter-spacing: -0.5px;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 16px;
  z-index: 1;
  margin-bottom: 80px;
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

const SecondaryButton = styled(Link)`
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: #fff;
  padding: 16px 32px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: background 0.2s;
  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
`;

const VisualContainer = styled.div`
  width: 100%;
  max-width: 1000px;
  height: 500px;
  background: radial-gradient(circle at center, #161b22 0%, #0d1117 100%);
  border: 1px solid #30363d;
  border-bottom: none;
  border-radius: 24px 24px 0 0;
  position: relative;
  overflow: hidden;
  box-shadow: 0 -20px 100px -20px rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Scanline = styled.div`
  position: absolute;
  left: 0;
  right: 0;
  height: 20%;
  background: linear-gradient(
    to bottom,
    transparent,
    rgba(88, 166, 255, 0.3),
    transparent
  );
  animation: ${scanlineAnim} 4s linear infinite;
  pointer-events: none;
  z-index: 10;
`;

// --- BENTO GRID WITH SPOTLIGHT ---
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

const FeatureSection = styled.section`
  padding: 100px 24px;
  max-width: 1200px;
  margin: 0 auto;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  @media (max-width: 1024px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
`;

const Card = styled.div<{ $colSpan?: number; $highlight?: string }>`
  background: rgba(13, 17, 23, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  grid-column: span ${(p) => p.$colSpan || 1};
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(12px);
  min-height: 320px;

  /* Spotlight Effect Logic */
  &::before {
    content: "";
    position: absolute;
    height: 100%;
    width: 100%;
    top: 0;
    left: 0;
    background: radial-gradient(
      800px circle at var(--mouse-x) var(--mouse-y),
      rgba(255, 255, 255, 0.06),
      transparent 40%
    );
    opacity: 0;
    transition: opacity 0.5s;
    z-index: 1;
    pointer-events: none;
  }

  &:hover::before {
    opacity: 1;
  }

  &:hover {
    border-color: ${(p) => p.$highlight || "#58a6ff"};
    transform: translateY(-2px);
    box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.5);
  }
`;

const CardContent = styled.div`
  position: relative;
  z-index: 2;
`;

const CardIcon = styled.div<{ $color: string }>`
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: ${(p) => `${p.$color}15`};
  color: ${(p) => p.$color};
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  border: 1px solid ${(p) => `${p.$color}30`};
`;

const CardTitle = styled.h3`
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
`;

const CardText = styled.p`
  color: #8b949e;
  font-size: 15px;
  line-height: 1.6;
  max-width: 90%;
`;

const SoftImage = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  animation: ${float} 6s ease-in-out infinite;
  mask-image: radial-gradient(circle at center, black 40%, transparent 100%);
  mix-blend-mode: lighten;
  opacity: 0.6;
  pointer-events: none;
`;

// --- PRICING ---
const PricingSection = styled.section`
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

// --- MAIN COMPONENT ---
export default function MarketingPage() {
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

  // Mouse Move for Spotlight
  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    const cards = document.querySelectorAll(".bento-card");
    cards.forEach((card) => {
      const htmlCard = card as HTMLElement;
      const rect = htmlCard.getBoundingClientRect();
      const x = e.clientX - rect.left;
      const y = e.clientY - rect.top;
      htmlCard.style.setProperty("--mouse-x", `${x}px`);
      htmlCard.style.setProperty("--mouse-y", `${y}px`);
    });
  };

  return (
    <Container onMouseMove={handleMouseMove}>
      {/* 1. HERO */}
      <HeroSection>
        <HeroGlow />
        <Pill>
          <Zap size={14} fill="currentColor" /> Sijil v1.0 Production Ready
        </Pill>

        <HeroTitle>
          <HeroRollingText />
        </HeroTitle>

        <Subtitle>
          Built by developers, for
          <HandDrawnHighlight> developers. </HandDrawnHighlight>
          <UPS>
            <span style={{ color: "#2ea043", fontWeight: 800 }}>104,000+</span>{" "}
            Logs/Second. Single Node.
          </UPS>
        </Subtitle>

        <ButtonGroup>
          <PrimaryButton href="/register">
            Start Ingesting <ArrowRight size={18} />
          </PrimaryButton>
          <SecondaryButton href="https://github.com/muhammad0320/sijil-core">
            View Source
          </SecondaryButton>
        </ButtonGroup>

        <VisualContainer>
          <Image
            src={server}
            alt="Engine"
            fill
            style={{
              objectFit: "cover",
              opacity: 0.9,
              mixBlendMode: "lighten",
            }}
          />
          <Scanline />

          <div
            style={{
              position: "absolute",
              inset: 0,
              background:
                "linear-gradient(to top, #050505 0%, transparent 50%)",
            }}
          />
        </VisualContainer>
      </HeroSection>

      {/* Becnchmark */}
      <BenchmarkSection />

      {/* 2. FEATURES (BENTO) */}
      <FeatureSection id="features">
        <SectionHeader>
          <h2>Built for the Paranoid</h2>
          <p>We assume everything will crash. That&apos;s why it works.</p>
        </SectionHeader>

        <Grid>
          {/* Feature 1 */}
          <Card $colSpan={2} $highlight="#58a6ff" className="bento-card">
            <CardContent>
              <CardIcon $color="#58a6ff">
                <Cpu size={24} />
              </CardIcon>
              <CardTitle>104k RPS Ingestion Engine</CardTitle>
              <CardText>
                Optimized Go engine with zero-allocation hot paths. We moved
                serialization outside the lock and use `goccy/go-json` for raw
                speed. Perfect for HFT spikes.
              </CardText>
            </CardContent>
            <div
              style={{
                position: "absolute",
                right: "-20px",
                bottom: "-40px",
                width: "300px",
                height: "300px",
                opacity: 0.4,
              }}
            >
              <SoftImage>
                <Image
                  src={hypercube}
                  alt="Hypercube"
                  fill
                  style={{ objectFit: "contain" }}
                />
              </SoftImage>
            </div>
          </Card>

          {/* Feature 2 */}
          <Card $highlight="#f1c40f" className="bento-card">
            <CardContent>
              <CardIcon $color="#f1c40f">
                <Shield size={24} />
              </CardIcon>
              <CardTitle>In-Memory Auth</CardTitle>
              <CardText>
                Project keys cached in RAM with aggressive TTLs. Nanosecond auth
                latency.
              </CardText>
            </CardContent>
            <div
              style={{
                position: "absolute",
                bottom: "-20px",
                right: "-20px",
                width: "160px",
                height: "160px",
                opacity: 0.3,
              }}
            >
              <SoftImage>
                <Image
                  src={cacheImg}
                  alt="Cache"
                  fill
                  style={{ objectFit: "contain" }}
                />
              </SoftImage>
            </div>
          </Card>

          {/* Feature 3 */}
          <Card $highlight="#2ecc71" className="bento-card">
            <CardContent>
              <CardIcon $color="#2ecc71">
                <Zap size={24} />
              </CardIcon>
              <CardTitle>Live Stream</CardTitle>
              <CardText>
                Real-time event loop. We optimized the broadcast hub to ignore
                empty rooms, saving CPU cycles.
              </CardText>
            </CardContent>
            <div
              style={{
                position: "absolute",
                bottom: "-20px",
                right: "-20px",
                width: "160px",
                height: "160px",
                opacity: 0.3,
              }}
            >
              <SoftImage>
                <Image
                  src={websocket}
                  alt="WS"
                  fill
                  style={{ objectFit: "contain" }}
                />
              </SoftImage>
            </div>
          </Card>

          {/* Feature 4 */}
          <Card $colSpan={2} $highlight="#a371f7" className="bento-card">
            <CardContent>
              <CardIcon $color="#a371f7">
                <Database size={24} />
              </CardIcon>
              <CardTitle>Write-Ahead Log (WAL)</CardTitle>
              <CardText>
                Buffered WAL ensures durability without killing IOPS. Recovery
                re-plays log segments on startup. Zero data loss.
              </CardText>
            </CardContent>
            <div
              style={{
                position: "absolute",
                right: "0",
                top: "0",
                height: "100%",
                width: "50%",
                opacity: 0.15,
              }}
            >
              <SoftImage>
                <Image
                  src={vault}
                  alt="WAL"
                  fill
                  style={{ objectFit: "contain", objectPosition: "right" }}
                />
              </SoftImage>
            </div>
          </Card>
        </Grid>
      </FeatureSection>

      {/* 3. PRICING */}
      <PricingSection id="pricing">
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
      </PricingSection>
    </Container>
  );
}
