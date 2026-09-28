export interface AgentEnvelope {
  schema: "kingdomsitter.agent-envelope/v0";
  task: string;
  analysis: {
    backend: string;
    variables: Array<{ name: string; type: string }>;
    calls: Array<{ name: string; effect: "read-or-output" | "potential-mutation" | "unknown"; count: number }>;
    effects: string[];
    raw_statements: number;
  };
  policy: {
    execution_authority: false;
    allowed_outputs: Array<"analysis" | "patch" | "typescript" | "explanation">;
  };
}

export type AgentResult =
  | { kind: "analysis"; markdown: string }
  | { kind: "typescript"; source: string }
  | { kind: "patch"; unifiedDiff: string };
