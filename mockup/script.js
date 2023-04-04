// Sample data representing the number of blocks mined by a player
const data = [
    {block: 'Stone', count: 150},
    {block: 'Dirt', count: 100},
    {block: 'Wood', count: 75},
    {block: 'Iron', count: 25},
    {block: 'Diamond', count: 10}
];

// Dimensions for the chart
const width = 280;
const height = 200;
const margin = {top: 20, right: 20, bottom: 30, left: 40};

// Create the SVG container for the chart
const svg = d3.select("#chart")
    .append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
    .append("g")
    .attr("transform", `translate(${margin.left}, ${margin.top})`);

// Set up scales for the chart
const x = d3.scaleBand()
    .domain(data.map(d => d.block))
    .range([0, width])
    .padding(0.1);
const y = d3.scaleLinear()
    .domain([0, d3.max(data, d => d.count)])
    .range([height, 0]);

// Create the bars for the chart
svg.selectAll(".bar")
    .data(data)
    .enter()
    .append("rect")
    .attr("class", "bar")
    .attr("x", d => x(d.block))
    .attr("width", x.bandwidth())
    .attr("y", d => y(d.count))
    .attr("height", d => height - y(d.count))
    .attr("fill", "#86795a");  // Minecraft grass block color

// Add the X-axis
svg.append("g")
    .attr("transform", `translate(0, ${height})`)
    .call(d3.axisBottom(x));

// Add the Y-axis
svg.append("g")
    .call(d3.axisLeft(y));

// Add chart title
svg.append("text")
    .attr("x", width / 2)
    .attr("y", -10)
    .attr("text-anchor", "middle")
    .text("Blocks Mined");
