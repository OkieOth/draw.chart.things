scriptPos=${0%/*}

cd $scriptPos/..

yacgVersion=6.10.0


# generate all the things from the business model
if ! docker run -u $(id -u ${USER}):$(id -g ${USER}) -v `pwd`:/project --rm -t \
    ghcr.io/okieoth/yacg:${yacgVersion} \
     --config /project/configs/codegen/config.json; then
    echo "error while run codegen, cancel"
    exit 1
fi

if ! docker run --rm -t \
    -u $(id -u ${USER}):$(id -g ${USER}) \
    -v $(pwd)/docs/models/puml:/puml \
    -v $(pwd)/docs/models/png:/output \
    -v $baseDir/docs:/docs \
    plantuml/plantuml \
    -tpng -O /output /puml; then
    echo "error while generate svg for: $fileToConvert"
    exit 1
fi

