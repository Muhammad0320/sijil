"use client";

import styled, { keyframes } from "styled-components";
import {
  ArrowRight,
  Database,
  Zap,
  Shield,
  Cpu,
  Activity,
  Check,
} from "lucide-react";
import Link from "next/link";
import Image from "next/image"; // Import Image for your assets
import hypercube from "../../../public/hypercube.png"; // Example image import
import valut from "../../../public/vault-green.png"; // Example image import
import serverExploded from "../../../public/server-exploded.png"; // Example image import

// --- ANIMATIONS ---
const scanline = keyframes`
  0% { transform: translateY(-100%); }
  100% { transform: translateY(100%); }
`;

// --- STYLES ---
const Container = styled.div`
  /* We don't need background color here, layout handles it */
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
`;

const HeroTitle = styled.h1`
  font-size: 84px;
  font-weight: 800;
  letter-spacing: -3px;
  line-height: 1.05;
  text-align: center;
  margin-bottom: 24px;
  z-index: 1;
  max-width: 90%; /* Prevent overflow */

  /* Gradient text... */
  background: linear-gradient(180deg, #fff 0%, #8b949e 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;

  span {
    color: #58a6ff;
    -webkit-text-fill-color: #58a6ff;
  }

  /* Tablet & Mobile Adjustment */
  @media (max-width: 1024px) {
    font-size: 56px; /* Smaller for tablets */
  }

  @media (max-width: 768px) {
    font-size: 42px; /* Smaller for mobile */
  }
`;

const HeroSection = styled.section`
  min-height: 90vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  /* padding-top: 40px; */
  overflow: hidden;

  /* The Grid Background */
  background-size: 50px 50px;
  background-image: linear-gradient(
      to right,
      rgba(255, 255, 255, 0.03) 1px,
      transparent 1px
    ),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  mask-image: radial-gradient(circle at center, black 40%, transparent 100%);
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
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  z-index: 0;
  filter: blur(80px);
`;

const Pill = styled.div`
  background: rgba(48, 54, 61, 0.5);
  border: 1px solid #30363d;
  padding: 6px 16px;
  border-radius: 99px;
  font-size: 12px;
  color: #58a6ff;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
  backdrop-filter: blur(4px);
  z-index: 1;
`;

const Subtitle = styled.p`
  font-size: 20px;
  color: #8b949e;
  text-align: center;
  max-width: 600px;
  margin-bottom: 48px;
  line-height: 1.6;
  z-index: 1;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 16px;
  z-index: 1;
  margin-bottom: 60px;
`;

const PrimaryButton = styled(Link)`
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  color: #000;
  padding: 14px 28px;
  border-radius: 8px;
  font-weight: 600;
  font-size: 15px;
  transition: transform 0.2s;

  &:hover {
    transform: translateY(-2px);
  }
`;

const SectionHeader = styled.div`
  margin-bottom: 64px;
  text-align: center;

  h2 {
    font-size: 48px;
    font-weight: 700;
    margin-bottom: 16px;
    letter-spacing: -1px;
  }

  p {
    color: #8b949e;
    font-size: 18px;
  }
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;

  /* Tablet: 2 Columns */
  @media (max-width: 1024px) {
    grid-template-columns: repeat(2, 1fr);
  }

  /* Mobile: 1 Column */
  @media (max-width: 640px) {
    grid-template-columns: 1fr;
  }
`;

const Card = styled.div<{ $colSpan?: number; $highlight?: string }>`
  background: rgba(22, 27, 34, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  grid-column: span ${(p) => p.$colSpan || 1};
  position: relative;
  overflow: hidden;
  transition: transform 0.3s ease, border-color 0.3s ease;

  /* The Spotlight Effect */
  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: radial-gradient(
      800px circle at var(--mouse-x, 50%) var(--mouse-y, 50%),
      rgba(255, 255, 255, 0.06),
      transparent 40%
    );
    z-index: 0;
    opacity: 0;
    transition: opacity 0.3s;
  }

  &:hover::before {
    opacity: 1;
  }

  &:hover {
    border-color: ${(p) => p.$highlight || "#58a6ff"};
    transform: translateY(-4px);
  }

  /* Reset spans on tablet/mobile so cards don't leave gaps */
  @media (max-width: 1024px) {
    grid-column: span 1 !important; /* Force single column width */
  }
`;

const CardContent = styled.div`
  position: relative;
  z-index: 1;
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
  font-weight: 600;
  margin-bottom: 12px;
  color: #fff;
`;

const CardText = styled.p`
  color: #8b949e;
  font-size: 15px;
  line-height: 1.6;
`;

const Scanline = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(
    to bottom,
    transparent,
    rgba(88, 166, 255, 0.1),
    transparent
  );
  animation: ${scanline} 3s linear infinite;
  pointer-events: none;
`;

// 4. Fix Section Padding
const BentoSection = styled.section`
  padding: 100px 24px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%; /* Ensure it doesn't overflow container */

  @media (max-width: 768px) {
    padding: 60px 16px;
  }
`;

// --- NEW STYLES FOR PERFORMANCE ---
const ComparisonSection = styled.section`
  padding: 100px 24px;
  max-width: 1000px;
  margin: 0 auto;
`;

const BarContainer = styled.div`
  margin-top: 48px;
  display: flex;
  flex-direction: column;
  gap: 24px;
`;

const BarGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const BarLabel = styled.div`
  display: flex;
  justify-content: space-between;
  color: #c9d1d9;
  font-size: 14px;
  font-weight: 500;
`;

const BarTrack = styled.div`
  width: 100%;
  height: 24px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
`;

const BarFill = styled.div<{ $width: string; $color: string }>`
  height: 100%;
  width: ${(p) => p.$width};
  background: ${(p) => p.$color};
  border-radius: 4px;
  position: relative;

  /* Moving scanline for the "fast" bar */
  &::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.4),
      transparent
    );
    transform: translateX(-100%);
    animation: ${(p) => (p.$width === "100%" ? "shimmer 2s infinite" : "none")};
  }

  @keyframes shimmer {
    100% {
      transform: translateX(100%);
    }
  }
`;

// --- NEW STYLES FOR PRICING ---
const PricingSection = styled.section`
  padding: 100px 24px;
  max-width: 1200px;
  margin: 0 auto 100px;
`;

const PricingGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  margin-top: 64px;

  @media (max-width: 860px) {
    grid-template-columns: 1fr;
  }
`;

const PriceCard = styled.div<{ $featured?: boolean }>`
  background: ${(p) =>
    p.$featured ? "rgba(88, 166, 255, 0.1)" : "rgba(22, 27, 34, 0.6)"};
  border: 1px solid ${(p) => (p.$featured ? "#58a6ff" : "#30363d")};
  border-radius: 16px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  position: relative;
  backdrop-filter: blur(12px);
  transition: transform 0.3s;

  &:hover {
    transform: translateY(-8px);
  }

  ${(p) =>
    p.$featured &&
    `
    &::before {
      content: 'MOST POPULAR';
      position: absolute;
      top: -12px;
      left: 50%;
      transform: translateX(-50%);
      background: #58a6ff;
      color: #000;
      font-size: 10px;
      font-weight: 800;
      padding: 4px 12px;
      border-radius: 20px;
    }
  `}
`;

const Price = styled.div`
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  margin: 16px 0;
  font-family: var(--font-geist-mono);

  span {
    font-size: 16px;
    color: #8b949e;
    font-weight: 400;
  }
`;

const FeatureList = styled.ul`
  list-style: none;
  margin: 24px 0 32px;
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

const FeatureItem = styled.li`
  display: flex;
  align-items: center;
  gap: 12px;
  color: #c9d1d9;
  font-size: 14px;
`;
// --- COMPONENT ---
export default function MarketingPage() {
  // FIXED: Strict Typing for Event
  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    // FIXED: Use querySelectorAll and proper casting
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
      <HeroSection>
        <HeroGlow />
        <Pill>
          <Zap size={12} fill="currentColor" /> v1.0.0 Stable Release
        </Pill>

        <HeroTitle>
          Ingest at <br />
          <span>Hyperspeed.</span>
        </HeroTitle>

        <Subtitle>
          The observability platform for the AI era. Built with Go, TimescaleDB,
          and a Write-Ahead Log that never loses a byte.
        </Subtitle>

        <ButtonGroup>
          <PrimaryButton href="/register">
            Get Started <ArrowRight size={16} />
          </PrimaryButton>
        </ButtonGroup>

        {/* HERO IMAGE CONTAINER */}
        <div
          style={{
            width: "100%",
            maxWidth: "1000px",
            height: "500px",
            background: "rgba(22,27,34,0.3)",
            border: "1px solid #30363d",
            borderBottom: "none",
            borderRadius: "24px 24px 0 0",
            marginTop: "40px",
            position: "relative",
            overflow: "hidden",
          }}
        >
          {/* Replace this src with your Exploded Server image */}
          <Image
            src={serverExploded}
            alt="Architecture"
            fill
            style={{ objectFit: "cover", opacity: 0.8 }}
          />
          <div
            style={{
              position: "absolute",
              inset: 0,
              background:
                "linear-gradient(to top, #050505 0%, transparent 50%)",
            }}
          />
          <Scanline />
        </div>
      </HeroSection>

      <BentoSection id="features">
        <SectionHeader>
          <h2>Engineered for Scale</h2>
          <p>We didn&apos;t just build a dashboard. We built a data engine.</p>
        </SectionHeader>

        <Grid>
          {/* Card 1: Hypertable */}
          <Card $colSpan={2} $highlight="#58a6ff" className="bento-card">
            <CardContent>
              <CardIcon $color="#58a6ff">
                <Database size={24} />
              </CardIcon>
              <CardTitle>TimescaleDB Hypertable</CardTitle>
              <CardText>
                Your logs are time-series data. Our engine automatically
                partitions data by time and project, enabling queries across
                billions of rows.
              </CardText>
            </CardContent>

            {/* Visual: Hypercube Image */}
            <div
              style={{
                position: "absolute",
                right: "-50px",
                bottom: "-50px",
                width: "300px",
                height: "300px",
                opacity: 0.5,
              }}
            >
              <Image
                src={hypercube}
                alt="Hypertable"
                fill
                style={{ objectFit: "contain" }}
              />
            </div>
          </Card>

          {/* Card 2: Security */}
          <Card $highlight="#2ecc71" className="bento-card">
            <CardContent>
              <CardIcon $color="#2ecc71">
                <Shield size={24} />
              </CardIcon>
              <CardTitle>Vault-Grade Auth</CardTitle>
              <CardText>
                API Keys are hashed using SHA-256 and salt. Granular RBAC
                ensures users only see what they own.
              </CardText>
            </CardContent>
          </Card>

          {/* Card 3: Realtime */}
          <Card $highlight="#f1c40f" className="bento-card">
            <CardContent>
              <CardIcon $color="#f1c40f">
                <Activity size={24} />
              </CardIcon>
              <CardTitle>Zero-Latency Stream</CardTitle>
              <CardText>
                Forget polling. Our WebSocket engine pushes logs directly from
                the ingestion pipeline instantly.
              </CardText>
            </CardContent>
          </Card>

          {/* Card 4: WAL Durability */}
          <Card $colSpan={2} $highlight="#a371f7" className="bento-card">
            <CardContent>
              <CardIcon $color="#a371f7">
                <Cpu size={24} />
              </CardIcon>
              <CardTitle>Powered by Write-Ahead Log</CardTitle>
              <CardText>
                We buffer every byte to a disk-based WAL before flushing. If the
                database crashes, your data survives.
              </CardText>
            </CardContent>
            {/* Visual: Vault Image */}
            <div
              style={{
                position: "absolute",
                right: "0",
                top: "0",
                width: "100%",
                height: "100%",
                opacity: 0.2,
                maskImage: "linear-gradient(to left, black, transparent)",
              }}
            >
              <Image
                src={valut}
                alt="WAL Durability"
                fill
                style={{ objectFit: "contain" }}
              />
            </div>
          </Card>
        </Grid>
      </BentoSection>

      <ComparisonSection id="performance">
        <SectionHeader>
          <h2>Benchmarks</h2>
          <p>Built for speed. We measure latency in microseconds.</p>
        </SectionHeader>

        <BarContainer>
          <BarGroup>
            <BarLabel>
              <span>LogEngine (Go + Hypertable)</span>
              <span style={{ color: "#2ecc71" }}>28ms (Ingest P99)</span>
            </BarLabel>
            <BarTrack>
              <BarFill $width="100%" $color="#2ecc71" />
            </BarTrack>
          </BarGroup>

          <BarGroup>
            <BarLabel>
              <span>Legacy Java Solution</span>
              <span>450ms</span>
            </BarLabel>
            <BarTrack>
              <BarFill $width="15%" $color="#30363d" />
            </BarTrack>
          </BarGroup>

          <BarGroup>
            <BarLabel>
              <span>Traditional SQL</span>
              <span>820ms</span>
            </BarLabel>
            <BarTrack>
              <BarFill $width="8%" $color="#30363d" />
            </BarTrack>
          </BarGroup>
        </BarContainer>
      </ComparisonSection>

      <PricingSection id="pricing">
        <SectionHeader>
          <h2>Transparent Pricing</h2>
          <p>Start small, scale infinitely.</p>
        </SectionHeader>

        <PricingGrid>
          {/* Free Plan */}
          <PriceCard>
            <CardTitle>Developer</CardTitle>
            <Price>
              $0<span>/mo</span>
            </Price>
            <CardText>For side projects and experiments.</CardText>
            <FeatureList>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> 3 Day Retention
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> 1 Member
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> 10k Logs/Day
              </FeatureItem>
            </FeatureList>
            <PrimaryButton
              href="/register"
              style={{
                width: "100%",
                justifyContent: "center",
                background: "rgba(255,255,255,0.1)",
                color: "#fff",
                border: "1px solid #30363d",
              }}
            >
              Start Free
            </PrimaryButton>
          </PriceCard>

          {/* Pro Plan */}
          <PriceCard $featured>
            <CardTitle>Pro Team</CardTitle>
            <Price>
              $29<span>/mo</span>
            </Price>
            <CardText>For startups shipping to production.</CardText>
            <FeatureList>
              <FeatureItem>
                <Check size={16} color="#58a6ff" /> 30 Day Retention
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#58a6ff" /> 10 Members
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#58a6ff" /> 1M Logs/Day
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#58a6ff" /> Email Alerts
              </FeatureItem>
            </FeatureList>
            <PrimaryButton
              href="/register"
              style={{ width: "100%", justifyContent: "center" }}
            >
              Get Pro
            </PrimaryButton>
          </PriceCard>

          {/* Ultra Plan */}
          <PriceCard>
            <CardTitle>Enterprise</CardTitle>
            <Price>Custom</Price>
            <CardText>For high-volume data compliance.</CardText>
            <FeatureList>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> 1 Year Retention
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> Unlimited Members
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> SSO & SAML
              </FeatureItem>
              <FeatureItem>
                <Check size={16} color="#2ecc71" /> Dedicated VPC
              </FeatureItem>
            </FeatureList>
            <PrimaryButton
              href="mailto:sales@logengine.com"
              style={{
                width: "100%",
                justifyContent: "center",
                background: "rgba(255,255,255,0.1)",
                color: "#fff",
                border: "1px solid #30363d",
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
