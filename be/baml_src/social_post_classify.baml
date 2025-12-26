// Social post classification with multimodal support (text + images)

// Input structure for a social post
class SocialPostInput {
  content string @description("The text content of the social post")
}

// Confidence levels for classification (using enum per BAML best practices)
enum ConfidenceLevel {
  HIGH @description("Very confident in this classification")
  MEDIUM @description("Moderately confident, some ambiguity present")
  LOW @description("Low confidence, classification is uncertain")
}

// Category types for social news posts
enum PostCategory {
  NEWS @description("Breaking news, current events, journalism")
  ENTERTAINMENT @description("Celebrity news, movies, music, pop culture")
  SPORTS @description("Sports news, game results, athlete updates")
  TECHNOLOGY @description("Tech news, gadgets, software, AI")
  BUSINESS @description("Business news, finance, markets, economy")
  POLITICS @description("Political news, government, elections")
  HEALTH @description("Health news, medical, wellness, fitness")
  SCIENCE @description("Scientific discoveries, research, space")
  LIFESTYLE @description("Food, travel, fashion, home")
  OPINION @description("Editorial, commentary, personal views")
  HUMOR @description("Jokes, memes, satirical content")
  ADVERTISEMENT @description("Promotional content, ads, sponsored posts")
  SPAM @description("Unwanted, irrelevant, or malicious content")
  OTHER @description("Content that doesn't fit other categories")
}

// Single classification result with category and confidence
class CategoryClassification {
  category PostCategory @description("The identified category for the post")
  confidence ConfidenceLevel @description("How confident the model is in this classification")
}

// Function to classify social posts
function ClassifySocialPost(post: string) -> CategoryClassification[] {
  client CustomGPT5Mini
  prompt #"
    Analyze the following social media post and classify it into relevant categories.
    Consider both the text content and any attached images to determine the categories.
    
    A post may belong to multiple categories. Return all applicable categories with your confidence level for each.
    Only include categories that genuinely apply to the content.

    {{ _.role("user") }}
    Post Content: {{ post }}


    {{ ctx.output_format }}
  "#
}

test TestName {
  functions [ClassifySocialPost]
  args {
    post {
        content #"
        hello world
      "#
    }
  }
}
