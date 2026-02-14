import sys
import requests


def create_svg(year=2025, round_num=1):
    url = f"http://localhost:8080/api/schedule/{year}/{round_num}/circuit"
    print(f"Fetching data from {url}...")

    try:
        response = requests.get(url)
        response.raise_for_status()
        data = response.json()
    except Exception as e:
        print(f"Error fetching data: {e}")
        sys.exit(1)

    layout = data.get('layout', [])
    if not layout:
        print("No layout data found in API response.")
        sys.exit(1)

    # Calculate bounds
    min_x = min(p['x'] for p in layout)
    max_x = max(p['x'] for p in layout)
    min_y = min(p['y'] for p in layout)
    max_y = max(p['y'] for p in layout)

    width = max_x - min_x
    height = max_y - min_y

    # Add some padding
    padding = max(width, height) * 0.1
    view_box = f"{min_x - padding} {min_y - padding} {width + 2*padding} {height + 2*padding}"

    # Create SVG path
    path_data = "M " + " L ".join(f"{p['x']},{p['y']}" for p in layout) + " Z"

    # Use circuit name in filename
    circuit_name = data.get(
        'circuit_name', 'circuit').replace(' ', '_').lower()
    output_file = f"{circuit_name}_layout.svg"

    # Create SVG content with a gradient
    svg_content = f"""<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg viewBox="{view_box}" xmlns="http://www.w3.org/2000/svg" style="background-color: #111;">
    <defs>
        <linearGradient id="trackGradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" style="stop-color:#ff4d4d;stop-opacity:1" />
            <stop offset="50%" style="stop-color:#ffcc00;stop-opacity:1" />
            <stop offset="100%" style="stop-color:#4dff4d;stop-opacity:1" />
        </linearGradient>
    </defs>
    <path d="{path_data}"
          stroke="url(#trackGradient)"
          stroke-width="{max(width, height) * 0.02}"
          fill="none"
          stroke-linecap="round"
          stroke-linejoin="round"/>

    <!-- Start/Finish Dot (approximate, first point) -->
    <circle cx="{layout[0]['x']}" cy="{layout[0]['y']}" r="{max(width, height) * 0.03}" fill="white" />
</svg>
"""

    with open(output_file, "w") as f:
        f.write(svg_content)

    print(f"SVG created as {output_file}")


if __name__ == "__main__":
    year_arg = int(sys.argv[1]) if len(sys.argv) > 1 else 2025
    round_arg = int(sys.argv[2]) if len(sys.argv) > 2 else 1
    create_svg(year_arg, round_arg)
