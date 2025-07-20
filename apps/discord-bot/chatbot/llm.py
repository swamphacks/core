from dotenv import load_dotenv
import os
from google import genai

load_dotenv()

client = genai.Client()

def llm_response(prompt):
    try:
        response = client.models.generate_content(
            model="gemini-1.5-flash",
            contents=prompt,
        )
        return response.text
    except Exception as e:
        return f"An error occurred: {str(e)}"
      
if __name__ == "__main__":
    prompt = input("Ask Gemini: ")
    print(llm_response(prompt))