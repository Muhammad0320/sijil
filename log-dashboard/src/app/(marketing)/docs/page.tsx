"use client";

import React, { useState, useEffect } from "react";
import styled, { keyframes, css } from "styled-components";
import { Copy, Check, Terminal, Key, Play, Shield } from "lucide-react";
import { HandDrawnHighlight } from "@/components/marketing/HeroText";

// --- ANIMATIONS ---
const fadeIn = keyframes`
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
`;

const blink = keyframes`
  50% { opacity: 0; }
`;

// --- STYLES ---
const Container = styled.div`
  min-height: 100vh;
  background: #050505;
  color: #c9d1d9;
  padding-top: 100px;
  padding-bottom: 100px;
  font-family: var(--font-geist-sans);
`;

const ContentWrapper = styled.div`
  max-width: 900px;
  margin: 0 auto;
  padding: 0 24px;
`;

const Header = styled.div`
  margin-bottom: 64px;
  text-align: center;
`;

const Title = styled.h1`
  font-size: 48px;
  font-weight: 800;
  color: #fff;
  margin-bottom: 16px;
  letter-spacing: -1px;
`;

const Subtitle = styled.p`
  font-size: 18px;
  color: #8b949e;
  max-width: 600px;
  margin: 0 auto;
`;

// --- INTERACTIVE API KEY INPUT ---
const ApiKeyCard = styled.div`
  background: rgba(13, 17, 23, 0.8);
  border: 1px solid #30363d;
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 64px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  backdrop-filter: blur(12px);
  position: sticky;
  top: 100px;
  z-index: 10;
  box-shadow: 0 20px 40px -10px rgba(0, 0, 0, 0.5);
`;

const InputRow = styled.div`
  display: flex;
  gap: 24px;
  align-items: flex-start;

  @media (max-width: 768px) {
    flex-direction: column;
    gap: 16px;
  }
`;

const InputGroup = styled.div`
  flex: 1;
  width: 100%;
`;

const Label = styled.div`
  font-size: 13px;
  font-weight: 600;
  color: #8b949e;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
`;

const InputWrapper = styled.div`
  position: relative;
`;

const Input = styled.input`
  width: 100%;
  background: #0d1117;
  border: 1px solid #30363d;
  color: #58a6ff;
  font-family: var(--font-geist-mono);
  padding: 10px 12px 10px 36px;
  border-radius: 8px;
  font-size: 13px;
  transition: all 0.2s;

  &:focus {
    outline: none;
    border-color: #58a6ff;
    box-shadow: 0 0 0 3px rgba(88, 166, 255, 0.1);
  }
  &::placeholder {
    color: #484f58;
  }
`;

// --- DOCS SECTION ---
const Section = styled.div`
  margin-bottom: 80px;
`;

const Step = styled.div`
  display: flex;
  gap: 32px;
  margin-bottom: 64px;

  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

const StepInfo = styled.div`
  flex: 1;
  max-width: 280px;
  padding-top: 12px;
`;

const StepNumber = styled.div`
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #161b22, #0d1117);
  border: 1px solid #30363d;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #fff;
  margin-bottom: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
`;

const StepTitle = styled.h3`
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
`;

const StepDesc = styled.p`
  font-size: 14px;
  color: #8b949e;
  line-height: 1.6;
`;

const StepCodeWrapper = styled.div`
  flex: 2;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #30363d;
  background: #0d1117;
  position: relative;

  /* The Resend-style Gradient */
  background: radial-gradient(
      circle at top right,
      rgba(88, 166, 255, 0.08),
      transparent 40%
    ),
    #0d1117;
`;

// --- TABS & CODE BLOCK ---
const TabHeader = styled.div`
  display: flex;
  border-bottom: 1px solid #30363d;
  background: rgba(22, 27, 34, 0.5);
  backdrop-filter: blur(4px);
`;

const TabBtn = styled.button<{ $active: boolean }>`
  padding: 12px 24px;
  background: ${(p) =>
    p.$active ? "rgba(88, 166, 255, 0.05)" : "transparent"};
  color: ${(p) => (p.$active ? "#58a6ff" : "#8b949e")};
  border: none;
  border-right: 1px solid #30363d;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    color: #fff;
  }
  ${(p) => p.$active && "box-shadow: inset 0 2px 0 #58a6ff;"}
`;

const CodeArea = styled.div`
  padding: 24px;
  position: relative;
  overflow-x: auto;
  font-family: var(--font-geist-mono);
  font-size: 13px;
  line-height: 1.7;
`;

// Syntax Highlighting Components
const K = styled.span`
  color: #ff7b72;
`; // Keyword
const S = styled.span`
  color: #a5d6ff;
`; // String
const F = styled.span`
  color: #d2a8ff;
`; // Function
const C = styled.span`
  color: #8b949e;
  font-style: italic;
`; // Comment
const P = styled.span`
  color: #79c0ff;
`; // Property
const V = styled.span`
  color: #c9d1d9;
`; // Variable

// --- LIVE DASHBOARD SIMULATION ---
const TerminalWindow = styled.div`
  background: #0d1117;
  padding: 20px;
  border-radius: 8px;
  font-family: var(--font-geist-mono);
  font-size: 12px;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const LogLine = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
  animation: ${fadeIn} 0.3s ease-out;
  border-bottom: 1px solid rgba(48, 54, 61, 0.4);
  padding-bottom: 8px;
  &:last-child {
    border-bottom: none;
  }
`;

const LogTime = styled.span`
  color: #8b949e;
`;
const LogLevel = styled.span<{ $level: string }>`
  font-weight: 700;
  color: ${(p) =>
    p.$level === "ERROR"
      ? "#ff7b72"
      : p.$level === "WARN"
      ? "#d29922"
      : "#2ea043"};
`;
const LogMsg = styled.span`
  color: #c9d1d9;
  flex: 1;
`;
const Cursor = styled.span`
  display: inline-block;
  width: 6px;
  height: 14px;
  background: #58a6ff;
  animation: ${blink} 1s infinite;
`;

const SimulateBtn = styled.button`
  position: absolute;
  bottom: 16px;
  right: 16px;
  background: #238636;
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s;
  &:hover {
    background: #2ea043;
    transform: translateY(-1px);
  }
  &:active {
    transform: translateY(0);
  }
`;

export default function DocsPage() {
  const [apiKey, setApiKey] = useState("");
  const [apiSecret, setApiSecret] = useState("");
  const [lang, setLang] = useState<"js" | "py" | "go">("js");
  const [copied, setCopied] = useState(false);
  const [logs, setLogs] = useState<any[]>([]);

  // Default placeholders
  const keyDisplay = apiKey || "pk_live_xxxxxxxxxxxxxxxx";
  const secretDisplay = apiSecret || "sk_live_xxxxxxxxxxxxxxxx";

  const copyToClipboard = () => {
    // Logic to extract text from current view (simplified)
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const simulateLog = () => {
    const newLog = {
      id: Date.now(),
      time: new Date().toLocaleTimeString(),
      level: Math.random() > 0.8 ? "ERROR" : "INFO",
      msg:
        Math.random() > 0.8
          ? "Payment gateway timeout"
          : "Payment processed successfully",
      service: "payment-service",
    };
    setLogs((prev) => [newLog, ...prev].slice(0, 5));
  };

  return (
    <Container>
      <ContentWrapper>
        <Header>
          <Title>Documentation</Title>
          <Subtitle>
            Integrate <HandDrawnHighlight>Sijil</HandDrawnHighlight> into your
            application in less than 2 minutes.
          </Subtitle>
        </Header>

        {/* INTERACTIVE INPUTS */}
        <ApiKeyCard>
          <div
            style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "center",
            }}
          >
            <h3
              style={{
                fontSize: 16,
                fontWeight: 700,
                color: "#fff",
                margin: 0,
              }}
            >
              Configure Your Snippets
            </h3>
            <div style={{ fontSize: 12, color: "#8b949e" }}>
              *Keys are local-only
            </div>
          </div>

          <InputRow>
            <InputGroup>
              <Label>
                <Key size={14} /> API Key (Public)
              </Label>
              <InputWrapper>
                <div
                  style={{
                    position: "absolute",
                    top: 11,
                    left: 12,
                    opacity: 0.5,
                  }}
                >
                  <Key size={14} />
                </div>
                <Input
                  placeholder="pk_live_..."
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                />
              </InputWrapper>
            </InputGroup>

            <InputGroup>
              <Label>
                <Shield size={14} /> API Secret (Private)
              </Label>
              <InputWrapper>
                <div
                  style={{
                    position: "absolute",
                    top: 11,
                    left: 12,
                    opacity: 0.5,
                  }}
                >
                  <Shield size={14} />
                </div>
                <Input
                  placeholder="sk_live_..."
                  type="password"
                  value={apiSecret}
                  onChange={(e) => setApiSecret(e.target.value)}
                />
              </InputWrapper>
            </InputGroup>
          </InputRow>
        </ApiKeyCard>

        {/* STEP 1: INSTALL */}
        <Section>
          <Step>
            <StepInfo>
              <StepNumber>1</StepNumber>
              <StepTitle>Install the SDK</StepTitle>
              <StepDesc>
                Install the Sijil client for your language. Zero dependencies.
              </StepDesc>
            </StepInfo>
            <StepCodeWrapper>
              <TabHeader>
                <TabBtn $active={lang === "js"} onClick={() => setLang("js")}>
                  npm
                </TabBtn>
                <TabBtn $active={lang === "py"} onClick={() => setLang("py")}>
                  pip
                </TabBtn>
                <TabBtn $active={lang === "go"} onClick={() => setLang("go")}>
                  go get
                </TabBtn>
              </TabHeader>
              <CodeArea>
                {lang === "js" && <div>npm install @sijil/node</div>}
                {lang === "py" && <div>pip install sijil</div>}
                {lang === "go" && <div>go get github.com/sijil/go-sdk</div>}
              </CodeArea>
            </StepCodeWrapper>
          </Step>

          {/* STEP 2: INITIALIZE (Using React Components for Coloring) */}
          <Step>
            <StepInfo>
              <StepNumber>2</StepNumber>
              <StepTitle>Send a Log</StepTitle>
              <StepDesc>
                Initialize the client and send a structured log.
              </StepDesc>
            </StepInfo>
            <StepCodeWrapper>
              <TabHeader>
                <TabBtn $active={lang === "js"} onClick={() => setLang("js")}>
                  Node.js
                </TabBtn>
                <TabBtn $active={lang === "py"} onClick={() => setLang("py")}>
                  Python
                </TabBtn>
                <TabBtn $active={lang === "go"} onClick={() => setLang("go")}>
                  Go
                </TabBtn>
              </TabHeader>
              <CodeArea>
                {lang === "js" && (
                  <>
                    <div>
                      <K>import</K> {"{"} <V>SijilLogger</V> {"}"} <K>from </K>
                      <S>
                        {'"'}@sijil/node{'"'}
                      </S>
                      ;
                    </div>
                    <br />
                    <div>
                      <K>const</K> <V>logger</V> = <K>new</K> <F>SijilLogger</F>
                      ({"{"}
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>apiKey:</P>
                      <S>
                        {'"'}
                        {keyDisplay}
                        {'"'}
                      </S>
                      ,
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>apiSecret:</P>
                      <S>
                        {'"'}
                        {secretDisplay}
                        {'"'}{" "}
                      </S>
                    </div>
                    <div>{"}"});</div>
                    <br />
                    <div>
                      <K>await</K> <V>logger</V>.<F>info</F>(
                      <S>
                        {'"'}Payment processed{'"'}
                      </S>
                      , {"{"}
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>amount:</P> <V>000</V>,
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>currency:</P>
                      <S>
                        {'"'}NGN{'"'}
                      </S>
                    </div>
                    <div>{"}"});</div>
                  </>
                )}

                {lang === "py" && (
                  <>
                    <div>
                      <K>from </K> sijil <K>import</K> SijilLogger
                    </div>
                    <br />
                    <div>
                      <V>logger</V> = <F>SijilLogger</F>(
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>api_key</P>=
                      <S>
                        {'"'} {keyDisplay} {'"'}
                      </S>
                      ,
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <P>api_secret</P>=
                      <S>
                        {'"'} {secretDisplay} {'"'}
                      </S>
                    </div>
                    <div>)</div>
                    <br />
                    <div>
                      <V>logger</V>.<F>info</F>(
                      <S>
                        {'"'} Payment processed {'"'}
                      </S>
                      , {"{"}
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <S>
                        {" "}
                        {'"'} amount {'"'}{" "}
                      </S>
                      : <V>5000</V>,
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <S>
                        {'"'} currency {'"'}
                      </S>
                      :
                      <S>
                        {'"'} NGN {'"'}
                      </S>
                    </div>
                    <div>{"}"})</div>
                  </>
                )}

                {lang === "go" && (
                  <>
                    <div>
                      <K>package</K> main
                    </div>
                    <div>
                      <K>import</K>
                      <S>
                        {'"'} github.com/sijil/go-sdk {'"'}
                      </S>
                    </div>
                    <br />
                    <div>
                      <K>func</K> <F>main</F>() {"{"}
                    </div>
                    <div style={{ paddingLeft: 20 }}>
                      <V>logger</V> := sijil.<F>New</F>(sijil.Config{"{"}
                    </div>
                    <div style={{ paddingLeft: 40 }}>
                      <P>APIKey:</P>
                      <S>
                        {'"'} {keyDisplay} {'"'}
                      </S>
                      ,
                    </div>
                    <div style={{ paddingLeft: 40 }}>
                      <P>APISecret:</P>
                      <S>
                        {'"'} {secretDisplay} {'"'}
                      </S>
                      ,
                    </div>
                    <div style={{ paddingLeft: 20 }}>{"}"})</div>
                    <br />
                    <div style={{ paddingLeft: 20 }}>
                      <V>logger</V>.<F>Info</F>(
                      <S>
                        {'"'} Payment processed {'"'}
                      </S>
                      , <K>map</K>[<K>string</K>]<K>interface</K>
                      {"{"}
                    </div>
                    <div style={{ paddingLeft: 40 }}>
                      <S>
                        {'"'} amount {'"'}
                      </S>
                      : <V>5000</V>,
                    </div>
                    <div style={{ paddingLeft: 40 }}>
                      <S>
                        {'"'} currency {'"'}
                      </S>
                      :
                      <S>
                        {'"'} NGN {'"'}
                      </S>
                      ,
                    </div>
                    <div style={{ paddingLeft: 20 }}>{"}"})</div>
                    <div>{"}"}</div>
                  </>
                )}
              </CodeArea>
            </StepCodeWrapper>
          </Step>

          {/* STEP 3: LIVE DASHBOARD SIMULATION */}
          <Step>
            <StepInfo>
              <StepNumber>3</StepNumber>
              <StepTitle>Watch it Live</StepTitle>
              <StepDesc>
                Go to your dashboard. Logs appear instantly via our WebSocket
                engine. Try sending a test log!
              </StepDesc>
            </StepInfo>
            <StepCodeWrapper style={{ border: "1px solid #30363d" }}>
              <TerminalWindow>
                {logs.length === 0 && (
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                      justifyContent: "center",
                      height: 160,
                      color: "#8b949e",
                      opacity: 0.6,
                    }}
                  >
                    <Terminal size={32} style={{ marginBottom: 12 }} />
                    <div>Listening for incoming logs...</div>
                  </div>
                )}
                {logs.map((log) => (
                  <LogLine key={log.id}>
                    <LogTime>{log.time}</LogTime>
                    <LogLevel $level={log.level}>{log.level}</LogLevel>
                    <LogMsg>
                      [{log.service}] {log.msg}
                    </LogMsg>
                  </LogLine>
                ))}
                <div style={{ marginTop: "auto", paddingTop: 8 }}>
                  <span style={{ color: "#2ea043" }}>➜</span>
                  <span style={{ color: "#58a6ff" }}>~</span> <Cursor />
                </div>
              </TerminalWindow>
              <SimulateBtn onClick={simulateLog}>
                <Play size={12} fill="currentColor" /> Simulate Log Event
              </SimulateBtn>
            </StepCodeWrapper>
          </Step>
        </Section>
      </ContentWrapper>
    </Container>
  );
}
