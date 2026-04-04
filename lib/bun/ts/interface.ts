import { sleep } from "bun";

function send(payload: Record<string, any>) {
  console.log(JSON.stringify(payload));
}

function sendInit() {
  send({
    data: "init",
  });
}

let i: number = 0;

export function getExampleData() {
  return {
    data: "Hello golang! " + i++,
    extra: i,
  };
}

sendInit();

while (i < 20) {
  send(getExampleData());

  await sleep(400);
}
