"use client";

import React, { useState, useEffect, useRef } from "react";
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

// Import your assets
import hypercube from "../../../public/hypercube.png";
import vault from "../../../public/vault-green.png";
import serverExploded from "../../../public/server-exploded.png";
import cacheImg from "../../../public/cache.png";
import websocket from "../../../public/websocket.png";
import HeroRollingText, {
  HandDrawnHighlight,
} from "@/components/marketing/HeroText";
import BenchmarkSection from "@/components/marketing/benchmark";

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

const scanlineAnim = keyframes`
  0% { top: -20%; opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { top: 100%; opacity: 0; }
`;

const revealAnim = keyframes`
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
`;

// --- GLOBAL STYLES ---
const Container = styled.div`
  min-height: 100vh;
  width: 100%;
  overflow-x: hidden;
  background: #050505;
  color: #c9d1d9;
  font-family: var(--font-geist-sans);
`;

// --- REVEAL COMPONENT (The Modern Scroll Effect) ---
const RevealWrapper = styled.div<{ $delay?: number; $isVisible: boolean }>`
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);
  transition-delay: ${(p) => p.$delay || 0}ms;

  ${(p) =>
    p.$isVisible &&
    css`
      opacity: 1;
      transform: translateY(0);
    `}
`;

function Reveal({
  children,
  delay = 0,
}: {
  children: React.ReactNode;
  delay?: number;
}) {
  const [isVisible, setIsVisible] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.disconnect(); // Only animate once
        }
      },
      { threshold: 0.1 }
    );

    if (ref.current) observer.observe(ref.current);
    return () => observer.disconnect();
  }, []);

  return (
    <RevealWrapper ref={ref} $isVisible={isVisible} $delay={delay}>
      {children}
    </RevealWrapper>
  );
}

// --- HERO SECTION ---
const HeroSection = styled.section`
  min-height: 95vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
  padding-top: 80px;

  /* Subtle Grid Background */
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
  top: 50%;
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
  font-weight: 500;
  color: #58a6ff;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
  backdrop-filter: blur(4px);
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

  @media (max-width: 768px) {
    font-size: 48px;
  }
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
  margin-bottom: 60px;
`;

const PrimaryButton = styled(Link)`
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  color: #000;
  padding: 14px 32px;
  border-radius: 8px;
  font-weight: 700;
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
  padding: 14px 32px;
  border-radius: 8px;
  font-weight: 600;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: background 0.2s;
  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
`;

// --- VISUALIZATION COMPONENTS ---
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
`;

// --- BENTO GRID SECTION ---
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

const FeatureGrid = styled.div`
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

// RENAMED TO AVOID CONFLICT
const FeatureCard = styled.div<{ $colSpan?: number; $highlight?: string }>`
  background: rgba(13, 17, 23, 0.6);
  border: 1px solid rgba(48, 54, 61, 0.6);
  border-radius: 24px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  grid-column: span ${(p) => p.$colSpan || 1};
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  backdrop-filter: blur(12px);
  min-height: 300px;

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

// --- PRICING SECTION ---
const PricingSection = styled.section`
  padding: 100px 24px;
  max-width: 1200px;
  margin: 0 auto;
`;

const ToggleContainer = styled.div`
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 48px;
`;

const ToggleBtn = styled.button<{ $active: boolean }>`
  padding: 8px 16px;
  border-radius: 8px;
  background: ${(p) => (p.$active ? "#238636" : "#21262d")};
  color: #fff;
  font-weight: 600;
  border: 1px solid #30363d;
  transition: all 0.2s;
`;

const PricingGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 32px;
  @media (max-width: 860px) {
    grid-template-columns: 1fr;
  }
`;

// RENAMED TO AVOID CONFLICT
const PriceCard = styled.div<{ $featured?: boolean; $ultra?: boolean }>`
  background: ${(p) => (p.$featured ? "rgba(88, 166, 255, 0.05)" : "#161b22")};
  border: 1px solid ${(p) => (p.$featured ? "#58a6ff" : "#30363d")};
  border-radius: 16px;
  padding: 32px;
  display: flex;
  flex-direction: column;
  position: relative;
  transition: transform 0.3s;

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
  margin: 24px 0;
  font-family: var(--font-geist-mono);
  span {
    font-size: 16px;
    color: #8b949e;
    font-weight: 400;
  }
`;

const FeatureList = styled.ul`
  list-style: none;
  margin: 0 0 32px 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const FeatureItem = styled.li`
  display: flex;
  align-items: center;
  gap: 12px;
  color: #c9d1d9;
  font-size: 14px;
`;

// --- UTILS ---
const Typewriter = ({ words }: { words: string[] }) => {
  const [index, setIndex] = useState(0);
  const [subIndex, setSubIndex] = useState(0);
  const [reverse, setReverse] = useState(false);
  const [blink, setBlink] = useState(true);

  useEffect(() => {
    const timeout = setTimeout(
      () => {
        if (subIndex === words[index].length + 1 && !reverse) {
          setReverse(true);
          return;
        }
        if (subIndex === 0 && reverse) {
          setReverse(false);
          setIndex((prev) => (prev + 1) % words.length);
          return;
        }
        setSubIndex((prev) => prev + (reverse ? -1 : 1));
      },
      reverse ? 75 : 150
    );
    return () => clearTimeout(timeout);
  }, [subIndex, index, reverse, words]);

  useEffect(() => {
    const timeout = setInterval(() => setBlink((prev) => !prev), 500);
    return () => clearInterval(timeout);
  }, []);

  return (
    <span style={{ color: "#58a6ff" }}>
      {words[index].substring(0, subIndex)}
      <span style={{ opacity: blink ? 1 : 0 }}>|</span>
    </span>
  );
};

// --- MAIN PAGE ---
export default function MarketingPage() {
  const [currency, setCurrency] = useState<"USD" | "NGN">("USD");
  const [isAnnual, setIsAnnual] = useState(true);

  // Pricing Logic
  const rates = {
    usd: { pro: 20, ultra: 100 },
    ngn: { pro: 12500, ultra: 95000 },
  };
  const currentRate = currency === "USD" ? rates.usd : rates.ngn;
  const symbol = currency === "USD" ? "$" : "₦";

  const getPrice = (base: number) => {
    if (!isAnnual) return base.toLocaleString();
    return ((base * 10) / 12).toLocaleString(undefined, {
      maximumFractionDigits: 0,
    });
  };

  return (
    <Container>
      {/* 1. HERO */}
      <HeroSection>
        <HeroGlow />
        <Reveal>
          <Pill>
            <Zap size={14} /> Sijil v1.0 Production Ready
          </Pill>
        </Reveal>

        <Reveal delay={100}>
          <HeroTitle>
            <HeroRollingText />
          </HeroTitle>
        </Reveal>

        <Reveal delay={200}>
          <Subtitle>
            Built by developers, for{" "}
            <HandDrawnHighlight> developers </HandDrawnHighlight>.
            <br />
            <span
              style={{ color: "white", fontWeight: "bold", display: "block" }}
            >
              100,000+ Logs/Seconds. Sigle Node.
            </span>
          </Subtitle>
        </Reveal>

        <Reveal delay={300}>
          <ButtonGroup>
            <PrimaryButton href="/register">
              Start Ingesting <ArrowRight size={18} />
            </PrimaryButton>
            <SecondaryButton href="https://github.com/muhammad0320/sijil-core">
              View Source
            </SecondaryButton>
          </ButtonGroup>
        </Reveal>

        <Reveal delay={400}>
          <VisualContainer>
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
          </VisualContainer>
        </Reveal>
      </HeroSection>

      {/* BenchMark Section */}
      <BenchmarkSection />

      {/* 2. FEATURES (BENTO GRID) */}
      <FeatureSection>
        <Reveal>
          <SectionHeader>
            <h2>Built for the Paranoid</h2>
            <p>We assume everything will crash. That&apos;s why it works!</p>
          </SectionHeader>
        </Reveal>

        <FeatureGrid>
          {/* Engine */}
          <Reveal delay={100}>
            <FeatureCard $colSpan={2} $highlight="#58a6ff">
              <CardContent>
                <CardIcon $color="#58a6ff">
                  <Cpu size={24} />
                </CardIcon>
                <CardTitle>Zero-Allocation Ingestion</CardTitle>
                <CardText>
                  Go engine optimized to reduce GC pauses. Workers drain
                  gracefully on shutdown.
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
            </FeatureCard>
          </Reveal>

          {/* Auth */}
          <Reveal delay={200}>
            <FeatureCard $highlight="#f1c40f">
              <CardContent>
                <CardIcon $color="#f1c40f">
                  <Shield size={24} />
                </CardIcon>
                <CardTitle>In-Memory Auth</CardTitle>
                <CardText>
                  Project keys cached in RAM. Nanosecond auth latency.
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
            </FeatureCard>
          </Reveal>

          {/* WebSockets */}
          <Reveal delay={300}>
            <FeatureCard $highlight="#2ecc71">
              <CardContent>
                <CardIcon $color="#2ecc71">
                  <Zap size={24} />
                </CardIcon>
                <CardTitle>WebSockets</CardTitle>
                <CardText>
                  Real-time event loop. No polling. See logs as they happen.
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
            </FeatureCard>
          </Reveal>

          {/* WAL (The one you lost!) */}
          <Reveal delay={400}>
            <FeatureCard $colSpan={2} $highlight="#a371f7">
              <CardContent>
                <CardIcon $color="#a371f7">
                  <Database size={24} />
                </CardIcon>
                <CardTitle>Write-Ahead Log (WAL)</CardTitle>
                <CardText>
                  Every byte buffered to disk. Crash recovery replays WAL on
                  startup. Zero data loss.
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
            </FeatureCard>
          </Reveal>
        </FeatureGrid>
      </FeatureSection>

      {/* 3. PRICING SECTION */}
      <PricingSection>
        <Reveal>
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
                borderRadius: 8,
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
                borderRadius: 8,
                border: "1px solid #30363d",
              }}
            >
              <ToggleBtn $active={!isAnnual} onClick={() => setIsAnnual(false)}>
                Monthly
              </ToggleBtn>
              <ToggleBtn $active={isAnnual} onClick={() => setIsAnnual(true)}>
                Yearly{" "}
                <span style={{ color: "#4ade80", fontSize: 10 }}>(-17%)</span>
              </ToggleBtn>
            </div>
          </ToggleContainer>
        </Reveal>

        <PricingGrid>
          {/* Hobby */}
          <Reveal delay={100}>
            <PriceCard>
              <CardTitle>Hobby</CardTitle>
              <PriceAmount>Free</PriceAmount>
              <FeatureList>
                <FeatureItem>
                  <Check size={16} color="#2ecc71" /> 10k Logs/Day
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#2ecc71" /> 3 Day Retention
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#2ecc71" /> Community Support
                </FeatureItem>
              </FeatureList>
              <PrimaryButton
                href="/register"
                style={{ justifyContent: "center", width: "100%" }}
              >
                Start Free
              </PrimaryButton>
            </PriceCard>
          </Reveal>

          {/* Pro */}
          <Reveal delay={200}>
            <PriceCard $featured>
              <div
                style={{
                  position: "absolute",
                  top: 0,
                  right: 0,
                  background: "#58a6ff",
                  color: "#000",
                  fontSize: 10,
                  fontWeight: 800,
                  padding: "4px 12px",
                  borderRadius: "0 16px 0 16px",
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
              <FeatureList>
                <FeatureItem>
                  <Check size={16} color="#58a6ff" /> 1M Logs/Day
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#58a6ff" /> 14 Day Retention
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#58a6ff" /> Email Support
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
                Upgrade Pro
              </PrimaryButton>
            </PriceCard>
          </Reveal>

          {/* Ultra */}
          <Reveal delay={300}>
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
              <FeatureList>
                <FeatureItem>
                  <Check size={16} color="#ffd700" /> 20M Logs/Day
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#ffd700" /> 30 Day Retention
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#ffd700" /> Cold Storage (S3)
                </FeatureItem>
                <FeatureItem>
                  <Check size={16} color="#ffd700" /> Priority Support
                </FeatureItem>
              </FeatureList>
              <PrimaryButton
                href="mailto:sales@sijil.io"
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
          </Reveal>
        </PricingGrid>
      </PricingSection>
    </Container>
  );
}
