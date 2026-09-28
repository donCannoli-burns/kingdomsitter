import type { AgentEnvelope } from "./agent-bridge.js";

export class SidecarClient {
  constructor(readonly baseUrl = "http://127.0.0.1:10423") {}

  async health(): Promise<Record<string, unknown>> {
    const response = await fetch(`${this.baseUrl}/health`);
    if (!response.ok) throw new Error(`kingdomsitter health failed: ${response.status}`);
    return (await response.json()) as Record<string, unknown>;
  }

  async analyzeAsh(source: string): Promise<AgentEnvelope> {
    const response = await fetch(`${this.baseUrl}/v0/agent`, {
      method: "POST",
      headers: { "content-type": "text/plain; charset=utf-8" },
      body: source,
    });
    if (!response.ok) throw new Error(await response.text());
    return (await response.json()) as AgentEnvelope;
  }

  async emitTypeScript(source: string): Promise<string> {
    const response = await fetch(`${this.baseUrl}/v0/emit-ts`, {
      method: "POST",
      headers: { "content-type": "text/plain; charset=utf-8" },
      body: source,
    });
    if (!response.ok) throw new Error(await response.text());
    return response.text();
  }
}
