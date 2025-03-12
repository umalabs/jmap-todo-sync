const JMAP_ENDPOINT = 'http://localhost:8080/jmap';

async function jmapRequest(methodCalls) {
    const response = await fetch(JMAP_ENDPOINT, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            methodCalls: methodCalls,
        }),
    });
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return await response.json();
}

async function getCapabilities() {
    const methodCalls = [
        ["Core/getCapabilities", {}, "r0"],
    ];
    return await jmapRequest(methodCalls);
}

async function getSession() {
    const methodCalls = [
        ["Core/getSession", {}, "r0"],
    ];
    return await jmapRequest(methodCalls);
}

async function todoQuery() {
    const methodCalls = [
        ["Todo/query", { "accountId": "primary" }, "r1"],
        ["Todo/get", { "accountId": "primary", "#ids": { "resultOf": "r1", "name": "Todo/query", "path": "/ids" } }, "r2"],
    ];
    return await jmapRequest(methodCalls);
}

async function todoSet(create, update, destroy) {
    let methodArgs = { "accountId": "primary" };
    if (create) methodArgs.create = create;
    if (update) methodArgs.update = update;
    if (destroy) methodArgs.destroy = destroy;

    const methodCalls = [
        ["Todo/set", methodArgs, "r1"],
    ];
    return await jmapRequest(methodCalls);
}


export { getCapabilities, getSession, todoQuery, todoSet };
