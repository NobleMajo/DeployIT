import { file, sleep } from "bun"

let msgPrefix: string = "Hello golang! "

const dataImportPath = "./data.ts"

if (await file(dataImportPath).exists()) {
  const dataImport = await import(dataImportPath)
  msgPrefix = dataImport.msgPrefix
  console.log("data.ts found, using import data prefix: ", msgPrefix)
} else {
  console.error("data.ts not found, test without import data")
}

function send(payload: Record<string, any>) {
  console.log(JSON.stringify(payload))
}

function sendInit() {
  send({
    data: "init",
  })
}

let i: number = 0

export function getExampleData() {
  return {
    data: msgPrefix + i,
    extra: i++,
  }
}

export async function readStdin() {
  const textDecoder = new TextDecoder()

  let bufferString = ""
  for await (const chunk of Bun.stdin.stream()) {
    bufferString += textDecoder.decode(chunk as Uint8Array, { stream: true })

    while (bufferString.includes("\n")) {
      const lineBreakIndex = bufferString.indexOf("\n")

      const line = bufferString.substring(0, lineBreakIndex)
      bufferString = bufferString.substring(lineBreakIndex + 1)

      try {
        const msg = JSON.parse(line)
        send({ from_bun: true, received: msg })
      } catch (err) {
        send({
          from_bun: true,
          error: "invalid_json",
          line: line,
          errorData: "" + err,
        })
      }
    }
  }
}

readStdin()

sendInit()

while (i < 20) {
  send(getExampleData())

  await sleep(400)
}
