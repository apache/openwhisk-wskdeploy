/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

function main(params) {
    let step = params.$step || 0
    delete params.$step
    package_name = "conductorPackage1"
    switch (step) {
        case 0: return { action: package_name+'/triple', params, state: { $step: 1 } }
        case 1: return { action: package_name+'/increment', params, state: { $step: 2 } }
        case 2: return { params }
    }
}
