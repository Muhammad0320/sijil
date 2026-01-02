"use client";

import styled, { keyframes } from "styled-components";
import { ArrowRight, Database, Zap, Shield, Cpu, Check } from "lucide-react";
import Link from "next/link";
import Image from "next/image"; // Import Image for your assets
import hypercube from "../../../public/hypercube.png"; // Example image import
import vault from "../../../public/vault-green.png"; // Example image import
import serverExploded from "../../../public/server-exploded.png"; // Example image import
import cacheImg from "../../../public/cache.png"; // Example image import
import websocket from "../../../public/websocket.png"; // Example image import

import { useEffect, useState } from "react";

// --- ANIMATIONS ---
const blink = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
`;

const float = keyframes`
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
`;

const scanline = keyframes`
  0% { transform: translateY(-100%); }
  100% { transform: translateY(100%); }
`;

// --- STYLES ---
const Container = styled.div`
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
  background: #050505; /* Fallback */
`;

const HeroTitle = styled.h1`
  font-size: 80px;
  font-weight: 800;
  letter-spacing: -3px;
  line-height: 1.1;
  text-align: center;
  margin-bottom: 24px;
  z-index: 1;
  max-width: 1100px;

  background: linear-gradient(180deg, #fff 0%, #8b949e 100%);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;

  /* Responsive Text */
  @media (max-width: 1024px) {
    font-size: 56px;
  }
  @media (max-width: 768px) {
    font-size: 40px;
  }
`;

const TypewriterCursor = styled.span`
  display: inline-block;
  width: 4px;
  height: 1em;
  background-color: #58a6ff;
  margin-left: 4px;
  vertical-align: middle;
  animation: ${blink} 1s step-end infinite;
`;

const Highlight = styled.span`
  color: #58a6ff;
  -webkit-text-fill-color: #58a6ff;
`;

const HeroSection = styled.section`
  min-height: 95vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  padding-top: 60px;
  overflow: hidden;

  /* Grid Background */
  background-size: 40px 40px;
  background-image: linear-gradient(
      to right,
      rgba(255, 255, 255, 0.03) 1px,
      transparent 1px
    ),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.03) 1px, transparent 1px);
  mask-image: radial-gradient(circle at center, black 30%, transparent 80%);
`;

const HeroGlow = styled.div`
  position: absolute;
  width: 800px;
  height: 800px;
  background: radial-gradient(
    circle,
    rgba(88, 166, 255, 0.12) 0%,
    transparent 70%
  );
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  z-index: 0;
  filter: blur(100px);
`;

const Pill = styled.div`
  background: rgba(88, 166, 255, 0.1);
  border: 1px solid rgba(88, 166, 255, 0.3);
  padding: 6px 16px;
  border-radius: 99px;
  font-size: 13px;
  font-weight: 500;
  color: #58a6ff;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
  backdrop-filter: blur(4px);
  box-shadow: 0 0 20px rgba(88, 166, 255, 0.1);
`;

const Subtitle = styled.p`
  font-size: 20px;
  color: #8b949e;
  text-align: center;
  max-width: 640px;
  margin-bottom: 40px;
  line-height: 1.6;
  z-index: 1;
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
  transition: all 0.2s;
  box-shadow: 0 0 20px rgba(255, 255, 255, 0.2);

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 0 30px rgba(255, 255, 255, 0.3);
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
  transition: all 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
`;

const SectionHeader = styled.div`
  margin-bottom: 64px;
  text-align: center;

  h2 {
    font-size: 48px;
    font-weight: 800;
    margin-bottom: 16px;
    letter-spacing: -1.5px;
    color: #fff;
  }

  p {
    color: #8b949e;
    font-size: 18px;
  }
`;

const BentoSection = styled.section`
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
  border: 1px solid rgba(48, 54, 61, 0.6);
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  grid-column: span ${(p) => p.$colSpan || 1};
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  backdrop-filter: blur(12px);

  /* Spotlight */
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
    box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.5);
  }
  @media (max-width: 1024px) {
    grid-column: span 1 !important;
  }
`;

const CardContent = styled.div`
  position: relative;
  z-index: 2; /* Ensure content sits above the faded image */
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
  box-shadow: 0 0 15px ${(p) => `${p.$color}20`};
`;

const CardTitle = styled.h3`
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 12px;
  color: #fff;
  letter-spacing: -0.5px;
`;

// FIX 4: Text constraint
const CardText = styled.p`
  color: #8b949e;
  font-size: 15px;
  line-height: 1.6;
  max-width: 80%; /* Prevents text from hitting the image */
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

const FloatingVisual = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  animation: ${float} 6s ease-in-out infinite;
  pointer-events: none;
`;

const scanlineAnim = keyframes`
  0% { top: -20%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
`;

// --- STYLES ---

// FIX 2: Scanline Implementation
const Scanline = styled.div`
  position: absolute;
  left: 0;
  right: 0;
  height: 20%; /* Height of the beam */
  background: linear-gradient(
    to bottom,
    transparent,
    rgba(88, 166, 255, 0.3),
    transparent
  );
  /* The animation now moves 'top' from -20% to 100% of the parent */
  animation: ${scanlineAnim} 4s linear infinite;
  pointer-events: none;
  z-index: 10;
`;

// FIX 3: Soft Edge Image Container
const SoftImage = styled.div`
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  animation: ${float} 6s ease-in-out infinite;

  /* THIS is the magic: Fades the image edges into the card background */
  mask-image: radial-gradient(circle at center, black 40%, transparent 100%);
  -webkit-mask-image: radial-gradient(
    circle at center,
    black 40%,
    transparent 100%
  );

  /* Blend mode helps dark images merge with dark cards */
  mix-blend-mode: lighten;
  opacity: 0.6; /* Subtle look */
  pointer-events: none;
`;

function Typewriter({ words }: { words: string[] }) {
  const [index, setIndex] = useState(0);
  const [subIndex, setSubIndex] = useState(0);
  const [reverse, setReverse] = useState(false);

  useEffect(() => {
    if (index >= words.length) return;
    if (subIndex === words[index].length + 1 && !reverse) {
      setTimeout(() => setReverse(true), 1500);
      return;
    }
    if (subIndex === 0 && reverse) {
      setTimeout(() => {
        setReverse(false);
        setIndex((prev) => (prev + 1) % words.length);
      }, 0);
      return;
    }
    const timeout = setTimeout(
      () => {
        setSubIndex((prev) => prev + (reverse ? -1 : 1));
      },
      reverse ? 50 : 100
    );
    return () => clearTimeout(timeout);
  }, [subIndex, index, reverse, words]);

  return (
    <Highlight>
      {words[index].substring(0, subIndex)}
      <TypewriterCursor />
    </Highlight>
  );
}
// --- MAIN PAGE ---
export default function MarketingPage() {
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
      <HeroSection>
        <HeroGlow />
        <Pill>
          <Zap size={14} fill="currentColor" /> Sijil v1.0 Production Ready
        </Pill>

        <HeroTitle>
          The Log Engine for <br />
          <Typewriter
            words={[
              "90k Logs/Seconds.",
              "Realtime Debugging.",
              "Crash-Proof WAL.",
            ]}
          />
        </HeroTitle>
        <Subtitle>
          Stop overpaying for Datadog. Sijil is the
          <Highlight> high-performance</Highlight>, self-hostable log engine
          built for startups and indie hackers.
          <br />
          <span
            style={{
              fontSize: "14px",
              opacity: 0.8,
              marginTop: "8px",
              display: "block",
            }}
          >
            Open Source • Single Binary • TimescaleDB Backed
          </span>
        </Subtitle>

        <ButtonGroup>
          <PrimaryButton href="/register">
            Start Ingesting <ArrowRight size={18} />
          </PrimaryButton>
          <SecondaryButton href="https://github.com/muhammad0320/sijil-core">
            View Source
          </SecondaryButton>
        </ButtonGroup>

        {/* HERO VISUALIZATION */}
        <div
          style={{
            width: "100%",
            maxWidth: "1000px",
            height: "550px",
            background:
              "radial-gradient(circle at center, #161b22 0%, #0d1117 100%)",
            border: "1px solid #30363d",
            borderBottom: "none",
            borderRadius: "24px 24px 0 0",
            marginTop: "20px",
            position: "relative",
            overflow: "hidden",
            boxShadow: "0 -20px 100px -20px rgba(0,0,0,0.8)",
          }}
        >
          <Image
            src={serverExploded}
            alt="Pipeline"
            fill
            style={{
              objectFit: "cover",
              opacity: 0.8,
              mixBlendMode: "lighten",
            }}
          />
          <Scanline />
          <div
            style={{
              position: "absolute",
              inset: 0,
              background:
                "linear-gradient(to top, #050505 5%, transparent 60%)",
            }}
          />
          <div
            style={{
              position: "absolute",
              top: 0,
              left: 0,
              right: 0,
              height: "2px",
              background: "#58a6ff",
              boxShadow: "0 0 20px #58a6ff",
            }}
          />
        </div>
      </HeroSection>

      <BentoSection id="features">
        <SectionHeader>
          <h2>Built for the Paranoid</h2>
          <p>We assume everything will crash. That&apos;s why it works.</p>
        </SectionHeader>

        <Grid>
          {/* Feature 1: The Engine */}
          <Card $colSpan={2} $highlight="#58a6ff" className="bento-card">
            <CardContent>
              <CardIcon $color="#58a6ff">
                <Cpu size={24} />
              </CardIcon>
              <CardTitle>Zero-Allocation Ingestion</CardTitle>
              <CardText>
                Our Go ingestion engine is optimized to reduce GC pauses.
                Workers drain gracefully on shutdown, ensuring no request is
                dropped even during deployments.
              </CardText>
            </CardContent>
            {/* FIX: Floating Hypercube */}
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
          {/* Feature 2: Auth (NOW WITH CACHE IMAGE) */}
          <Card $highlight="#f1c40f" className="bento-card">
            <CardContent>
              <CardIcon $color="#f1c40f">
                <Shield size={24} />
              </CardIcon>
              <CardTitle>In-Memory Auth</CardTitle>
              <CardText>
                We don&apos;t hit the DB for every log. Project keys are cached
                in RAM with aggressive TTLs, cutting auth latency to
                nanoseconds.
              </CardText>
            </CardContent>
            {/* FIX: New Cache Image */}
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
                  alt="Cache Shield"
                  fill
                  style={{ objectFit: "contain" }}
                />
              </SoftImage>
            </div>
          </Card>

          {/* Feature 3: Live Stream */}
          <Card $highlight="#2ecc71" className="bento-card">
            <CardContent>
              <CardIcon $color="#2ecc71">
                <Zap size={24} />
              </CardIcon>
              <CardTitle>WebSockets</CardTitle>
              <CardText>
                Real-time event loop; See logs as they happen. Our event loop
                pipes ingestion directly to connected clients. No polling. No
                lag.
              </CardText>
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
                    alt="digital connection"
                    fill
                    style={{ objectFit: "contain" }}
                  />
                </SoftImage>
              </div>
            </CardContent>
          </Card>

          {/* Feature 4: WAL */}
          <Card $colSpan={2} $highlight="#a371f7" className="bento-card">
            <CardContent>
              <CardIcon $color="#a371f7">
                <Database size={24} />
              </CardIcon>
              <CardTitle>Write-Ahead Log (WAL)</CardTitle>
              <CardText>
                Every byte is buffered to disk before processing. If the server
                crashes, our recovery protocol replays the WAL on startup. Zero
                data loss.
              </CardText>
            </CardContent>
            {/* FIX: Floating Vault */}
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
