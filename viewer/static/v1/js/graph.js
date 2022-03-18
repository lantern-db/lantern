import * as d3 from "https://cdn.skypack.dev/d3@7";

export class Graph {
    constructor(data, param) {
        this.data = data;
        this.param = param;
    }

    start() {
        this.simulation = d3.forceSimulation(this.data.nodes)
            .force("charge", d3.forceManyBody().strength(-100))
            .force("link", d3.forceLink(this.data.links))
            .force("center", d3.forceCenter(0, 0))
            .force("collide", d3.forceCollide(this.param.r))
            .alpha(1.5)
            .on("tick", ticked);

        this.svg = d3.select("#canvas").append("svg")
            .attr("id", "graph")
            .attr("class", "w-100 h-100 mx-auto")
            .attr("viewBox", [-this.param.width / 2, -this.param.height / 2, this.param.width, this.param.height]);

        let links = this.svg.append("g")
            .attr("stroke", "dimgray")
            .attr("stroke-opacity", "0.2")
            .attr("stroke-width", "1.5")
            .attr("stroke-linecap", "round")
            .selectAll("line")
            .data(this.data.links)
            .join("line");

        this.nodes = this.svg.append("g")
            .selectAll("g")
            .data(this.data.nodes)
            .join("g")
            .on("click", (event, d) => {
                window.location.href = `/v1/?seed=${d.key}&step=${this.param.step}&k=${this.param.k}`;
            })
            .call(this.drag(this.simulation));

        const circles = this.nodes
            .append("circle")
            .attr("r", this.param.r);

        circles
            .filter((d, i) => d.key === this.param.seed)
            .attr("fill", "orange");

        const labels = this.nodes
            .append("text")
            .text(d => d.value);

        labels
            .filter((d, i) => d.key === this.param.seed)
            .attr("fill", "orange");

        function ticked() {
            links
                .attr("x1", d => d.source.x)
                .attr("y1", d => d.source.y)
                .attr("x2", d => d.target.x)
                .attr("y2", d => d.target.y);

            circles
                .attr("cx", d => d.x)
                .attr("cy", d => d.y);

            labels
                .attr("x", d => d.x)
                .attr("y", d => d.y);
        }
    }



    drag(simulation) {
        function dragstarted(event) {
            if (!event.active) simulation.alphaTarget(0.3).restart();
            event.subject.fx = event.subject.x;
            event.subject.fy = event.subject.y;
        }

        function dragged(event) {
            event.subject.fx = event.x;
            event.subject.fy = event.y;
        }

        function dragended(event) {
            if (!event.active) simulation.alphaTarget(0);
            event.subject.fx = null;
            event.subject.fy = null;
        }

        return d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended);
    }
}


export function transform(graph) {
    const nodes = graph.vertices.map((v, i) => {
        return {
            index: i,
            key: v.key,
            value: v["string"] === undefined ? `[${v.key}]` : v.string
        }
    });
    return {
        nodes: nodes,
        links: graph.edges.map((e, i) => {
            return {
                source: nodes.find(({key: key}) => key === e.tail),
                target: nodes.find(({key: key}) => key === e.head),
            }
        })
    }
}
