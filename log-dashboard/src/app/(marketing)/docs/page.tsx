"use client";

import React, { useState } from "react";
import styled from "styled-components";
import { Copy, Check, Terminal, Key } from "lucide-react";
import { HandDrawnHighlight } from "@/components/marketing/HeroText";

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
  background: rgba(13, 17, 23, 0.6);
  border: 1px solid #30363d;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 48px;
  display: flex;
  align-items: center;
  gap: 24px;
  backdrop-filter: blur(12px);
  position: sticky;
  top: 100px; /* Sticks below navbar */
  z-index: 10;
  box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.5);

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
`;

const Label = styled.div`
  font-weight: 600;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 140px;
`;

const InputWrapper = styled.div`
  flex: 1;
  position: relative;
`;

const Input = styled.input`
  width: 100%;
  background: #0d1117;
  border: 1px solid #30363d;
  color: #58a6ff;
  font-family: var(--font-geist-mono);
  padding: 10px 12px 10px 40px;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;

  &:focus {
    outline: none;
    border-color: #58a6ff;
  }
  &::placeholder {
    color: #484f58;
  }
`;

// --- DOCS SECTION ---
const Section = styled.div`
  margin-bottom: 80px;
`;

const SectionTitle = styled.h2`
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 12px;
`;

const Step = styled.div`
  display: flex;
  gap: 24px;
  margin-bottom: 48px;

  @media (max-width: 768px) {
    flex-direction: column;
  }
`;

const StepInfo = styled.div`
  flex: 1;
  max-width: 300px;
`;

const StepNumber = styled.div`
  width: 32px;
  height: 32px;
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  color: #8b949e;
  margin-bottom: 16px;
`;

const StepTitle = styled.h3`
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8px;
`;

const StepDesc = styled.p`
  font-size: 14px;
  color: #8b949e;
  line-height: 1.6;
`;

const StepCode = styled.div`
  flex: 2;
  background: #0d1117;
  border: 1px solid #30363d;
  border-radius: 12px;
  overflow: hidden;
`;

// --- TABS & CODE BLOCK ---
const TabHeader = styled.div`
  display: flex;
  border-bottom: 1px solid #30363d;
  background: #161b22;
`;

const TabBtn = styled.button<{ $active: boolean }>`
  padding: 10px 20px;
  background: ${(p) => (p.$active ? "#0d1117" : "transparent")};
  color: ${(p) => (p.$active ? "#58a6ff" : "#8b949e")};
  border: none;
  border-right: 1px solid #30363d;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;

  ${(p) => p.$active && "border-top: 2px solid #58a6ff; margin-top: -1px;"}
`;

const CodeArea = styled.div`
  padding: 24px;
  position: relative;
  overflow-x: auto;
`;

const Pre = styled.pre`
  margin: 0;
  font-family: var(--font-geist-mono);
  font-size: 13px;
  line-height: 1.6;
  color: #c9d1d9;
`;

const Keyword = styled.span`
  color: #ff7b72;
`;
const String = styled.span`
  color: #a5d6ff;
`;
const Func = styled.span`
  color: #d2a8ff;
`;
const Comment = styled.span`
  color: #8b949e;
  font-style: italic;
`;

// --- COMPONENT ---

export default function DocsPage() {
  const [apiKey, setApiKey] = useState("");
  const [lang, setLang] = useState<"js" | "py" | "go">("js");
  const [copied, setCopied] = useState(false);

  // Helper to inject the key or show placeholder
  const keyDisplay = apiKey || "pk_live_xxxxxxxxxxxxxxxx";
  const secretDisplay = apiKey
    ? apiKey.replace("pk_", "sk_")
    : "sk_live_xxxxxxxxxxxxxxxx"; // Just a visual mock

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const getCode = () => {
    if (lang === "js")
      return `
<Keyword>import</Keyword> { SijilLogger } <Keyword>from</Keyword> <String>"@sijil/node"</String>;

<Keyword>const</Keyword> logger = <Keyword>new</Keyword> <Func>SijilLogger</Func>({
  apiKey: <String>"${keyDisplay}"</String>,
  // In production, use environment variables!
});

<Keyword>await</Keyword> logger.<Func>info</Func>(<String>"Payment processed"</String>, {
  amount: 5000,
  currency: <String>"NGN"</String>,
  user_id: <String>"user_123"</String>
});`;

    if (lang === "py")
      return `
<Keyword>from</Keyword> sijil <Keyword>import</Keyword> SijilLogger

logger = <Func>SijilLogger</Func>(
    api_key=<String>"${keyDisplay}"</String>
)

<Comment># Send a structured log</Comment>
logger.<Func>info</Func>(<String>"Payment processed"</String>, {
    <String>"amount"</String>: 5000,
    <String>"currency"</String>: <String>"NGN"</String>
})`;

    if (lang === "go")
      return `
<Keyword>package</Keyword> main

<Keyword>import</Keyword> <String>"github.com/sijil/go-sdk"</String>

<Keyword>func</Keyword> <Func>main</Func>() {
    logger := sijil.<Func>New</Func>(sijil.Config{
        APIKey: <String>"${keyDisplay}"</String>,
    })

    logger.<Func>Info</Func>(<String>"Payment processed"</String>, <Keyword>map</Keyword>[<Keyword>string</Keyword>]<Keyword>interface</Keyword>{
        <String>"amount"</String>: 5000,
        <String>"currency"</String>: <String>"NGN"</String>,
    })
}`;
    return "";
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

        {/* INTERACTIVE INPUT */}
        <ApiKeyCard>
          <Label>
            <Key size={16} /> API Key
          </Label>
          <InputWrapper>
            <Input
              placeholder="Paste your pk_live_... key here"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
            />
            <div
              style={{ position: "absolute", top: 11, left: 12, opacity: 0.5 }}
            >
              <Key size={14} />
            </div>
          </InputWrapper>
          <div style={{ fontSize: 13, color: "#8b949e", maxWidth: 200 }}>
            *We don&apos;t store this. It just updates the snippets below.
          </div>
        </ApiKeyCard>

        {/* STEP 1: INSTALL */}
        <Section>
          <Step>
            <StepInfo>
              <StepNumber>1</StepNumber>
              <StepTitle>Install the SDK</StepTitle>
              <StepDesc>
                Install the Sijil client for your language of choice. Our SDKs
                are zero-dependency and lightweight.
              </StepDesc>
            </StepInfo>
            <StepCode>
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
                <Pre>
                  {lang === "js" && "npm install @sijil/node"}
                  {lang === "py" && "pip install sijil"}
                  {lang === "go" && "go get github.com/sijil/go-sdk"}
                </Pre>
              </CodeArea>
            </StepCode>
          </Step>

          {/* STEP 2: INITIALIZE */}
          <Step>
            <StepInfo>
              <StepNumber>2</StepNumber>
              <StepTitle>Send your first log</StepTitle>
              <StepDesc>
                Initialize the client and send a structured log. We
                automatically capture timestamps and service metadata.
              </StepDesc>
            </StepInfo>
            <StepCode>
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
                <button
                  onClick={() =>
                    copyToClipboard(getCode().replace(/<[^>]+>/g, ""))
                  }
                  style={{
                    position: "absolute",
                    top: 16,
                    right: 16,
                    background: "#21262d",
                    border: "1px solid #30363d",
                    padding: 6,
                    borderRadius: 6,
                    cursor: "pointer",
                    color: copied ? "#2ecc71" : "#8b949e",
                  }}
                >
                  {copied ? <Check size={14} /> : <Copy size={14} />}
                </button>
                <Pre dangerouslySetInnerHTML={{ __html: getCode() }} />
              </CodeArea>
            </StepCode>
          </Step>

          {/* STEP 3: DASHBOARD */}
          <Step>
            <StepInfo>
              <StepNumber>3</StepNumber>
              <StepTitle>Check the Dashboard</StepTitle>
              <StepDesc>
                Go to the{" "}
                <a href="/dashboard" style={{ color: "#58a6ff" }}>
                  Live Dashboard
                </a>
                . You should see your log appear instantly via WebSocket.
              </StepDesc>
            </StepInfo>
            <StepCode>
              <div
                style={{ padding: 40, textAlign: "center", color: "#8b949e" }}
              >
                <Terminal
                  size={48}
                  style={{ marginBottom: 16, opacity: 0.5 }}
                />
                <div>Waiting for incoming logs...</div>
              </div>
            </StepCode>
          </Step>
        </Section>
      </ContentWrapper>
    </Container>
  );
}
