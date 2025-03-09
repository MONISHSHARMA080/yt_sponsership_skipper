<script>
  
  
  let position = 0;
  let title = "Background Paths";

  const paths = Array.from({ length: 36 }, (_, i) => ({
    id: i,
    d: `M-${380 - i * 5 + position} -${189 + i * 6}C-${380 - i * 5 + position} -${189 + i * 6} -${312 - i * 5 + position} ${216 - i * 6} ${152 - i * 5 + position} ${343 - i * 6}C${616 - i * 5 + position} ${470 - i * 6} ${684 - i * 5 + position} ${875 - i * 6} ${684 - i * 5 + position} ${875 - i * 6}`,
    color: `rgba(15,23,42,${0.1 + i * 0.03})`,
    width: 0.5 + i * 0.03,
  }));

  const words = title.split(" ");
</script>

<main class="relative w-full h-screen flex flex-col items-center justify-center overflow-hidden">
  {#each paths as path}
    <motion.path
      d={path.d}
      stroke={path.color}
      stroke-width={path.width}
      fill="none"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 2, delay: path.id * 0.05 }}
    />
  {/each}

  <div class="text-center z-10 relative">
    {#each words as word, wordIndex}
      <motion.div
        class="inline-block text-4xl font-bold"
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: wordIndex * 0.2 }}>
        {#each word.split("") as letter, letterIndex}
          <motion.span
            class="inline-block"
            initial={{ opacity: 0, y: 5 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: (wordIndex + letterIndex) * 0.05 }}>
            {letter}
          </motion.span>
        {/each}
        &nbsp;
      </motion.div>
    {/each}
  </div>

  <motion.button
    class="mt-6 bg-blue-600 text-white px-6 py-3 rounded-lg shadow-lg"
    whileHover={{ scale: 1.05 }}
    whileTap={{ scale: 0.95 }}>
    Discover Excellence →
  </motion.button>
</main>

<style>
  main {
    background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
    color: white;
  }
</style>
