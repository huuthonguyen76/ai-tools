import { GoogleGenAI } from "@google/genai";
import { ContextualizedResult } from "../types";

const API_KEY = process.env.API_KEY || '';

export const generateLinkContext = async (url: string): Promise<ContextualizedResult> => {
  if (!API_KEY) {
    throw new Error("API Key is missing. Please ensure process.env.API_KEY is set.");
  }

  const ai = new GoogleGenAI({ apiKey: API_KEY });

  const prompt = `
    Analyze the following URL: ${url}

    Task: Create a "Contextualized Smart Link" that explains the content clearly in English.

    1. **Analyze**: Use Google Search to read the current content of the page.
    2. **Highlight**: Extract a UNIQUE, VERBATIM short phrase (5-10 words) from the main body text. It must be an exact match to create a working Text Fragment link. Do not change a single character.
    3. **Label**: Create a short, punchy, descriptive English title (3-6 words) that summarizes the specific topic. This will be converted into a URL slug (e.g., "React Hooks Introduction" -> "#react-hooks-introduction").
    4. **Context**: Write a clear, single-sentence summary in English explaining exactly what this link is about.
    5. **Tags**: 3 relevant English keywords.

    Format your response strictly as follows:
    ---TITLE---
    (The Page Title)
    ---HIGHLIGHT_PHRASE---
    (The verbatim text fragment)
    ---DESCRIPTIVE_LABEL---
    (The English Label)
    ---SUMMARY---
    (The English Context Summary)
    ---TAGS---
    (Tag1, Tag2, Tag3)
  `;

  try {
    const response = await ai.models.generateContent({
      model: 'gemini-2.5-flash',
      contents: prompt,
      config: {
        tools: [{ googleSearch: {} }],
        temperature: 0.3, 
      },
    });

    const text = response.text;
    
    if (!text) {
        throw new Error("No content generated.");
    }

    // Helper to robustly extract sections
    const getSection = (startMarker: string, endMarker: string) => {
        const regex = new RegExp(`${startMarker}([\\s\\S]*?)${endMarker}`);
        const match = text.match(regex);
        return match ? match[1].trim() : "";
    };

    const highlightPhrase = getSection('---HIGHLIGHT_PHRASE---', '---DESCRIPTIVE_LABEL---');
    // Clean up highlight phrase if it contains quotes which might break URL encoding later
    const cleanHighlight = highlightPhrase.replace(/["']/g, "");

    const tagsText = text.match(/---TAGS---([\s\S]*?)$/);
    const suggestedTags = tagsText ? tagsText[1].split(',').map(t => t.trim()) : [];

    // Extract Grounding Metadata for sources
    const sources = response.candidates?.[0]?.groundingMetadata?.groundingChunks
      ?.map((chunk: any) => chunk.web ? { title: chunk.web.title, uri: chunk.web.uri } : null)
      .filter((item: any) => item !== null) as { title: string; uri: string }[] || [];

    return {
      originalUrl: url,
      title: getSection('---TITLE---', '---HIGHLIGHT_PHRASE---') || "External Link",
      highlightPhrase: cleanHighlight,
      descriptiveLabel: getSection('---DESCRIPTIVE_LABEL---', '---SUMMARY---') || "Visit Link",
      summary: getSection('---SUMMARY---', '---TAGS---') || "Content details unavailable.",
      suggestedTags: suggestedTags,
      sources: sources,
    };

  } catch (error) {
    console.error("Gemini API Error:", error);
    throw error;
  }
};