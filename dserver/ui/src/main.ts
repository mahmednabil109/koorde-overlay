import { ACTION_GETBLOCKS, ACTION_INITBRAOD } from "./actions";
import { RenderBuilder } from "./renderBuilder";
import "./style.css";

let body = document.querySelector(".body") || { innerHTML: "" };
var renderer = new RenderBuilder(".body", 500, 12);

var nodes = new Map<string, [string, string]>();
var up = false;

var ws = new WebSocket("ws://localhost:8282/ws");
ws.addEventListener("open", (event) => {
  console.log(event);
  up = true;
});

ws.addEventListener("message", (event) => {
  let data = JSON.parse(event.data);

  if (data.blocks?.length) {
    // display the blocks
    displayBlocks(data.blocks);
  } else {
    // register new node and update
    nodes.set(data.id, [data.successor, data.d]);
    update();
  }
});

function update() {
  let n = [],
    p: [number, number][] = [];
  for (let node of nodes.keys()) n.push(parseInt(`0x${node.substring(0, 3)}`));
  for (const [key, val] of nodes) {
    for (let v of val)
      p.push([
        parseInt(`0x${key.substring(0, 3)}`),
        parseInt(`0x${v.substring(0, 3)}`),
      ]);
  }

  n = n.filter((v) => !isNaN(v));
  p = p.filter(([p1, p2]) => !isNaN(p1) && !isNaN(p2));

  body.innerHTML = ``;
  renderer.set_peers(n).set_path(p).render();
}

// handle actions
var selected = null;
var selectedID = null;

var broadcastBtn = document.querySelector("button#broadcast");
var getBlocksBtn = document.querySelector("button#get-blocks");

window.handleNodeClick = (e: number) => {
  let id = [...nodes.keys()][e];
  selected = document.querySelector(".selected-node .id");
  selected.textContent = id;
  selectedID = id;

  // if (up)
  //   ws.send(
  //     JSON.stringify({
  //       id: id,
  //       action: ACTION_GETBLOCKS,
  //     })
  //   );
};

broadcastBtn?.addEventListener("click", () => {
  if (!up || !selectedID) return;

  var block = document.querySelector("textarea#broadcast").value;

  ws.send(
    JSON.stringify({
      id: selectedID,
      action: ACTION_INITBRAOD,
      block: block,
    })
  );
});

getBlocksBtn?.addEventListener("click", () => {
  if (!up || !selectedID) return;

  ws.send(
    JSON.stringify({
      id: selectedID,
      action: ACTION_GETBLOCKS,
    })
  );
});

function displayBlocks(blocks) {
  let blocksContainer = document.querySelector(".blocks");
  let elements = "";

  for (let block of blocks) {
    elements += `
    <div class="block">${block}</div>
    `;
  }
  blocksContainer.innerHTML = elements;
}
