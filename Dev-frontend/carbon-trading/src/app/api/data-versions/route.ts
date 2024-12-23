// app/api/data-versions/route.ts
import { NextResponse } from "next/server";

export async function GET() {
  try {
    // Replace with your actual API endpoint
    const apiUrl = `${process.env.API_BASE_URL}/data-versions`;
    const apiKey = `${process.env.API_KEY}/data-versions`;
    const response = await fetch(apiUrl, {
      headers: {
        Authorization: `Bearer ${apiKey}`,
      },
    });

    const data = await response.json();
    return NextResponse.json(data);
  } catch {
    return NextResponse.json(
      { error: "Failed to fetch data versions" },
      { status: 500 }
    );
  }
}
