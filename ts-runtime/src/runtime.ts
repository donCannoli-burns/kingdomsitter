import type { AgentEnvelope } from "./agent-bridge.js";

/**
 * kingdomsitter deliberately does not hide an execution channel inside the
 * agent API. The runtime receives structured proposals; a separate host layer
 * decides whether/how anything reaches KoLmafia.
 */
export interface KingdomSitterRuntime {
  analyzeAsh(source: string): Promise<AgentEnvelope>;
  emitTypeScript(source: string): Promise<string>;
}
