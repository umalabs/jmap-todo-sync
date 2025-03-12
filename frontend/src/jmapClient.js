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
        // **Modified: Get IDs from Todo/query and pass them directly to Todo/get**
        ["Todo/get", { "accountId": "primary", "ids": [] }, "r2"], // Initially send empty ids, we'll populate later
    ];
    const response = await jmapRequest(methodCalls);
    if (response.methodResponses && response.methodResponses.length > 0) {
        const queryResponse = response.methodResponses.find(resp => resp[0] === 'Todo/query');
        if (queryResponse && queryResponse[1] && queryResponse[1].ids) {
            const ids = queryResponse[1].ids;
            // Now, update the 'Todo/get' method call in the response to include the ids
            const getMethodCall = response.methodResponses.find(resp => resp[0] === 'Todo/get');
            if (getMethodCall && getMethodCall[1]) {
                getMethodCall[1].ids = ids; // Set the 'ids' directly in the Todo/get response arguments
            }
        }
        return response;
    }
    return response;
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
