/**
 * Return a simple string to
 * confirm this function has been visited.
 *
 * @param visited the visited function list
 */
function main(params) {
    functionID = params.functionID || 'X'
    if (params.visited == null) {
        params.visited = 'function'+functionID;
    } else {
        params.visited = params.visited + ', function'+functionID;
    }
    return {"visited":params.visited};
}
