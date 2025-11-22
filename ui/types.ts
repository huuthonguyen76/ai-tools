export interface ContextualizedResult {
  originalUrl: string;
  title: string;
  highlightPhrase: string;
  descriptiveLabel: string;
  summary: string;
  suggestedTags: string[];
  sources?: { title: string; uri: string }[];
}

export interface GenerationState {
  isLoading: boolean;
  error: string | null;
  result: ContextualizedResult | null;
}