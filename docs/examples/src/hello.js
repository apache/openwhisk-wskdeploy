function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    return { payload:  msg };
}
