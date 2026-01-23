#!/usr/bin/env python3
"""
Test script to verify Discord bot can access production database via API.
This tests the same endpoint that assign_hacker_roles uses.
"""

import os
import sys
import asyncio
import aiohttp
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

API_URL = os.getenv("API_URL", "https://api.swamphacks.com")
SESSION_COOKIE = os.getenv("SESSION_COOKIE")
EVENT_ID = os.getenv("EVENT_ID")

async def test_api_connection():
    """Test basic API connectivity"""
    print(f"🔍 Testing API connection to: {API_URL}")
    
    async with aiohttp.ClientSession() as session:
        try:
            async with session.get(f"{API_URL}/ping") as response:
                if response.status == 200:
                    text = await response.text()
                    print(f"✅ API is reachable: {text.strip()}")
                    return True
                else:
                    print(f"❌ API ping failed with status: {response.status}")
                    return False
        except Exception as e:
            print(f"❌ Failed to connect to API: {e}")
            return False

async def test_attendees_endpoint():
    """Test the attendees endpoint that assign_hacker_roles uses"""
    if not SESSION_COOKIE:
        print("❌ SESSION_COOKIE is not set in .env file")
        return False
    
    if not EVENT_ID:
        print("❌ EVENT_ID is not set in .env file")
        return False
    
    print(f"\n🔍 Testing attendees endpoint:")
    print(f"   API URL: {API_URL}")
    print(f"   Event ID: {EVENT_ID}")
    print(f"   Session Cookie: {SESSION_COOKIE[:20]}...")
    
    endpoint = f"{API_URL}/discord/event/{EVENT_ID}/attendees"
    
    async with aiohttp.ClientSession() as session:
        headers = {"Cookie": f"sh_session_id={SESSION_COOKIE}"}
        
        try:
            async with session.get(endpoint, headers=headers) as response:
                status = response.status
                print(f"\n📊 Response Status: {status}")
                
                if status == 200:
                    data = await response.json()
                    attendee_count = len(data)
                    print(f"✅ Successfully retrieved {attendee_count} attendees from production database!")
                    
                    if attendee_count > 0:
                        print(f"\n📋 Sample attendee data:")
                        sample = data[0]
                        print(f"   - Discord ID: {sample.get('discord_id', 'N/A')}")
                        print(f"   - User ID: {sample.get('user_id', 'N/A')}")
                        print(f"   - Name: {sample.get('name', 'N/A')}")
                        print(f"   - Email: {sample.get('email', 'N/A')}")
                    
                    return True
                    
                elif status == 401:
                    print("❌ Authentication failed - SESSION_COOKIE may be invalid or expired")
                    text = await response.text()
                    print(f"   Response: {text}")
                    return False
                    
                elif status == 404:
                    print("⚠️  Event not found or no attendees with Discord IDs")
                    return True  # This is still a valid response
                    
                else:
                    text = await response.text()
                    print(f"❌ Request failed with status {status}")
                    print(f"   Response: {text}")
                    return False
                    
        except aiohttp.ClientError as e:
            print(f"❌ Network error: {e}")
            return False
        except Exception as e:
            print(f"❌ Unexpected error: {e}")
            import traceback
            traceback.print_exc()
            return False

async def main():
    print("=" * 60)
    print("Discord Bot Production Database Connection Test")
    print("=" * 60)
    
    # Test 1: Basic API connectivity
    api_ok = await test_api_connection()
    if not api_ok:
        print("\n❌ Cannot proceed - API is not reachable")
        sys.exit(1)
    
    # Test 2: Attendees endpoint (the one assign_hacker_roles uses)
    db_ok = await test_attendees_endpoint()
    
    print("\n" + "=" * 60)
    if db_ok:
        print("✅ All tests passed! Your bot can access the production database.")
        print("   The assign_hacker_roles command should work correctly.")
    else:
        print("❌ Tests failed. Please check:")
        print("   1. SESSION_COOKIE is valid and not expired")
        print("   2. EVENT_ID is correct")
        print("   3. Your session has proper permissions")
        sys.exit(1)
    print("=" * 60)

if __name__ == "__main__":
    asyncio.run(main())
